package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

type Config struct {
	Environment string
	ServiceName string

	KafkaUrl string

	PostgresHost     string
	PostgresPort     int
	PostgresUser     string
	PostgresPassword string
	PostgresDatabase string

	MinioAccessKeyID string
	MinioSecretKey   string
	MinioEndpoint    string

	LogLevel string

	HttpHost string
	HttpPort string
}

func Load() Config {
	envFileName := cast.ToString(getOrReturnDefault("ENV_FILE_PATH", "./app/.env"))

	if err := godotenv.Load(envFileName); err != nil {
		fmt.Println("No .env file found")
	}
	config := Config{}

	config.Environment = cast.ToString(getOrReturnDefault("ENVIRONMENT", "develop"))
	config.LogLevel = cast.ToString(getOrReturnDefault("LOG_LEVEL", "info"))
	config.ServiceName = cast.ToString(getOrReturnDefault("SERVICE_NAME", ""))

	config.MinioEndpoint = cast.ToString(getOrReturnDefault("MINIO_ENDPOINT", "dev.cdn.7i.uz"))
	config.MinioAccessKeyID = cast.ToString(getOrReturnDefault("MINIO_ACCESS_KEY_ID", "seeRah2mraiL3eehOiTi6eewUux6zohd"))
	config.MinioSecretKey = cast.ToString(getOrReturnDefault("MINIO_SECRET_KEY_ID", "keigie9Arae4AeShpaeg6Cheahd3CohY"))

	config.PostgresHost = cast.ToString(getOrReturnDefault("POSTGRES_HOST", "localhost"))
	config.PostgresPort = cast.ToInt(getOrReturnDefault("POSTGRES_PORT", 5433))
	config.PostgresUser = cast.ToString(getOrReturnDefault("POSTGRES_USER", "postgres"))
	config.PostgresPassword = cast.ToString(getOrReturnDefault("POSTGRES_PASSWORD", "root"))
	config.PostgresDatabase = cast.ToString(getOrReturnDefault("POSTGRES_DATABASE", "invan_corporate_service"))

	config.KafkaUrl = cast.ToString(getOrReturnDefault("KAFKA_URL", "localhost:9092"))

	config.HttpHost = cast.ToString(getOrReturnDefault("LISTEN_HOST", "localhost"))
	config.HttpPort = cast.ToString(getOrReturnDefault("GRPC_PORT", ":8003"))

	return config
}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	_, exists := os.LookupEnv(key)

	if exists {
		return os.Getenv(key)
	}

	return defaultValue
}
