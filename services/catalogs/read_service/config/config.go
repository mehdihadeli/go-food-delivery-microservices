package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/constants"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/elasticsearch"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/eventstroredb"
	kafkaClient "github.com/mehdihadeli/store-golang-microservice-sample/pkg/kafka"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mongodb"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/probes"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/redis"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "", "catalogs read microservice config path")
}

type Config struct {
	DeliveryType     string                         `mapstructure:"deliveryType"`
	ServiceName      string                         `mapstructure:"serviceName"`
	Logger           *logger.Config                 `mapstructure:"logger"`
	KafkaTopics      KafkaTopics                    `mapstructure:"kafkaTopics"`
	GRPC             GRPC                           `mapstructure:"grpc"`
	Http             Http                           `mapstructure:"http"`
	Context          Context                        `mapstructure:"context"`
	Redis            *redis.Config                  `mapstructure:"redis"`
	Rabbitmq         *rabbitmq.RabbitMQConfig       `mapstructure:"rabbitmq"`
	Kafka            *kafkaClient.Config            `mapstructure:"kafka"`
	Probes           probes.Config                  `mapstructure:"probes"`
	Jaeger           *tracing.Config                `mapstructure:"jaeger"`
	EventStoreConfig eventstroredb.EventStoreConfig `mapstructure:"eventStoreConfig"`
	Subscriptions    Subscriptions                  `mapstructure:"subscriptions"`
	Elastic          elasticsearch.Config           `mapstructure:"elastic"`
	ElasticIndexes   ElasticIndexes                 `mapstructure:"elasticIndexes"`
	Mongo            *mongodb.Config                `mapstructure:"mongo"`
	MongoCollections MongoCollections               `mapstructure:"mongoCollections"`
}

type Context struct {
	Timeout int `mapstructure:"timeout"`
}

type GRPC struct {
	Port        string `mapstructure:"port"`
	Development bool   `mapstructure:"development"`
}

type Http struct {
	Port                string   `mapstructure:"port" validate:"required"`
	Development         bool     `mapstructure:"development"`
	BasePath            string   `mapstructure:"basePath" validate:"required"`
	ProductsPath        string   `mapstructure:"productsPath" validate:"required"`
	DebugErrorsResponse bool     `mapstructure:"debugErrorsResponse"`
	IgnoreLogUrls       []string `mapstructure:"ignoreLogUrls"`
	Timeout             int      `mapstructure:"timeout"`
	Host                string   `mapstructure:"host"`
}

type MongoCollections struct {
	Products string `mapstructure:"products" validate:"required"`
}

type Subscriptions struct {
	PoolSize                   int    `mapstructure:"poolSize" validate:"required,gte=0"`
	OrderPrefix                string `mapstructure:"orderPrefix" validate:"required,gte=0"`
	MongoProjectionGroupName   string `mapstructure:"mongoProjectionGroupName" validate:"required,gte=0"`
	ElasticProjectionGroupName string `mapstructure:"elasticProjectionGroupName" validate:"required,gte=0"`
}

type ElasticIndexes struct {
	Orders string `mapstructure:"orders" validate:"required"`
}

type KafkaTopics struct {
	ProductCreate  kafkaClient.TopicConfig `mapstructure:"productCreate"`
	ProductCreated kafkaClient.TopicConfig `mapstructure:"productCreated"`
	ProductUpdate  kafkaClient.TopicConfig `mapstructure:"productUpdate"`
	ProductUpdated kafkaClient.TopicConfig `mapstructure:"productUpdated"`
	ProductDelete  kafkaClient.TopicConfig `mapstructure:"productDelete"`
	ProductDeleted kafkaClient.TopicConfig `mapstructure:"productDeleted"`
}

func InitConfig(env string) (*Config, error) {
	if configPath == "" {
		configPathFromEnv := os.Getenv(constants.ConfigPath)
		if configPathFromEnv != "" {
			configPath = configPathFromEnv
		} else {
			configPath = "./config"
		}
	}

	cfg := &Config{}

	//https://github.com/spf13/viper/issues/390#issuecomment-718756752
	viper.SetConfigName(fmt.Sprintf("config.%s", env))
	viper.AddConfigPath(configPath)
	viper.SetConfigType(constants.Yaml)

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "viper.ReadInConfig")
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, errors.Wrap(err, "viper.Unmarshal")
	}

	grpcPort := os.Getenv(constants.GrpcPort)
	if grpcPort != "" {
		cfg.GRPC.Port = grpcPort
	}

	mongoURI := os.Getenv(constants.MongoDbURI)
	if mongoURI != "" {
		//cfg.Mongo.URI = "mongodb://host.docker.internal:27017"
		cfg.Mongo.URI = mongoURI
	}

	redisAddr := os.Getenv(constants.RedisAddr)
	if redisAddr != "" {
		cfg.Redis.Addr = redisAddr
	}

	jaegerAddr := os.Getenv(constants.JaegerHostPort)
	if jaegerAddr != "" {
		cfg.Jaeger.HostPort = jaegerAddr
	}

	kafkaBrokers := os.Getenv(constants.KafkaBrokers)
	if kafkaBrokers != "" {
		cfg.Kafka.Brokers = []string{kafkaBrokers}
	}

	return cfg, nil
}
