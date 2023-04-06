package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"

	"genproto/corporate_service"

	"github.com/Invan2/invan_corporate_service/config"
	"github.com/Invan2/invan_corporate_service/events"
	"github.com/Invan2/invan_corporate_service/pkg/logger"
	"github.com/Invan2/invan_corporate_service/services/listeners"
	"github.com/Invan2/invan_corporate_service/storage"
	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
)

func main() {

	cfg := config.Load()
	log := logger.New(cfg.LogLevel, cfg.ServiceName)
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	log.Info("config", logger.Any("config", cfg), logger.Any("env", os.Environ()))

	postgresURL := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDatabase,
	)

	psqlConn, err := sqlx.Connect("postgres", postgresURL)
	if err != nil {
		log.Error("postgres", logger.Error(err))
		return
	}

	defer psqlConn.Close()

	strg := storage.NewStoragePg(log, psqlConn, cfg)

	conf := kafka.ConfigMap{
		"bootstrap.servers":                     cfg.KafkaUrl,
		"group.id":                              config.ConsumerGroupID,
		"auto.offset.reset":                     "earliest",
		"go.events.channel.size":                1000000,
		"socket.keepalive.enable":               true,
		"metadata.max.age.ms":                   900000,
		"metadata.request.timeout.ms":           30000,
		"retries":                               1000000,
		"message.timeout.ms":                    300000,
		"socket.timeout.ms":                     30000,
		"max.in.flight.requests.per.connection": 5,
		"heartbeat.interval.ms":                 3000,
		"enable.idempotence":                    true,
	}

	log.Info("kafka config", logger.Any("config", conf))

	producer, err := kafka.NewProducer(&conf)
	if err != nil {
		log.Error("error while creating producer")
		return
	}

	consumer, err := kafka.NewConsumer(&conf)
	if err != nil {
		log.Error("error while creating consumer", logger.Error(err))
		return
	}

	pubSub, err := events.NewPubSubServer(log, producer, consumer, strg)
	if err != nil {
		log.Error("pub sub", logger.Error(err))
		return
	}

	server := grpc.NewServer()
	corporate_service.RegisterCorporateServiceServer(server, listeners.NewCorporateService(log, pubSub, strg))

	lis, err := net.Listen("tcp", fmt.Sprintf("%s%s", cfg.HttpHost, cfg.HttpPort))
	if err != nil {
		log.Error("http", logger.Error(err))
		return
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-c
		fmt.Println("Gracefully shutting down...")
		server.GracefulStop()
		if err := pubSub.Shutdown(); err != nil {
			log.Error("error while shutdown pub sub server", logger.Error(err))
			return
		}
	}()

	go func() {
		if err := pubSub.Run(ctx); err != nil {
			panic(err)
		}
	}()

	if err := server.Serve(lis); err != nil {
		log.Fatal("serve", logger.Error(err))
		return
	}

}
