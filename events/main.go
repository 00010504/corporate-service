package events

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/Invan2/invan_corporate_service/events/handlers"
	"github.com/Invan2/invan_corporate_service/events/topics"
	"github.com/Invan2/invan_corporate_service/pkg/logger"
	"github.com/Invan2/invan_corporate_service/storage"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type pubSubServer struct {
	log           logger.Logger
	producer      *kafka.Producer
	consumerGroup *kafka.Consumer
	eventHandlers map[string]EventHandler
	strgPG        storage.StoragePgI
	ctx           context.Context
	wg            *sync.WaitGroup
}

type PubSubServer interface {
	Run(ctx context.Context) error
	Push(topic string, value interface{}) error
	AddConsumer(topic string, handler EventHandler)

	Shutdown() error
}

type EventHandler func(ctx context.Context, event *kafka.Message) error

func NewPubSubServer(log logger.Logger, producer *kafka.Producer, consumer *kafka.Consumer, strgPG storage.StoragePgI) (PubSubServer, error) {
	return &pubSubServer{
		log:           log,
		producer:      producer,
		consumerGroup: consumer,
		eventHandlers: make(map[string]EventHandler),
		strgPG:        strgPG,
		wg:            &sync.WaitGroup{},
	}, nil
}

func (p *pubSubServer) Push(topic string, value interface{}) error {

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	deliveryChan := make(chan kafka.Event, 10000)
	err = p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: data,
	}, deliveryChan)
	if err != nil {
		return err
	}

	return nil

}

func (p *pubSubServer) AddConsumer(topic string, handler EventHandler) {
	_, ok := p.eventHandlers[topic]
	if ok {
		panic(fmt.Sprintf("%v topic already exits", topic))
	}

	p.eventHandlers[topic] = handler
}

func (p *pubSubServer) registerConsumers() {
	handlerV1 := handlers.NewHandler(p.log, p.strgPG)

	p.AddConsumer(topics.UserUpsertTopic, handlerV1.Upsert)

}

func (p *pubSubServer) Run(ctx context.Context) error {
	p.ctx = ctx

	p.registerConsumers()

	go func() {
		if err := p.producerLogger(ctx); err != nil {
			p.log.Error("error while logging producer events", logger.Error(err))
		}
	}()

	topics := make([]string, 0)

	for topic := range p.eventHandlers {
		topics = append(topics, topic)
	}

	if len(topics) == 0 {
		return errors.New("no topics")
	}

	if err := p.consumerGroup.SubscribeTopics(topics, nil); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:

			e := p.consumerGroup.Poll(100)

			switch e := e.(type) {
			case *kafka.Message:
				handler, ok := p.eventHandlers[*e.TopicPartition.Topic]
				if !ok {
					p.log.Error("handler not fount", logger.String("topic", *e.TopicPartition.Topic))
					continue
				}

				p.wg.Add(1)

				go func() {

					defer p.wg.Done()

					if err := handler(ctx, e); err != nil {
						p.log.Error("error while handling event", logger.Error(err))
						return
					}
				}()
			case kafka.Error:
				p.log.Error("error while consuming", logger.Error(e))
			default:
			}

		}
	}

}

func (p *pubSubServer) Shutdown() error {
	select {
	case <-p.ctx.Done():
		p.log.Warn("terminating: context cancelled")
	default:
	}

	if err := p.consumerGroup.Close(); err != nil {
		return err
	}

	p.wg.Wait()

	return nil
}

func (p *pubSubServer) producerLogger(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			e := <-p.producer.Events()
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					p.log.Error("Failed to deliver message:", logger.Any("topic", ev.TopicPartition))
				}
			}
		}
	}
}
