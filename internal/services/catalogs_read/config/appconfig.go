package config

import (
	"os"
	"strings"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/constants"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/grpc"
	customEcho "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/custom_echo"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/metrics"
	rabbitmqconfig "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/config"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/elasticsearch"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/eventstroredb"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/mongodb"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/redis"
)

type AppConfig struct {
	DeliveryType      string                         `mapstructure:"deliveryType"     env:"DeliveryType"`
	ServiceName       string                         `mapstructure:"serviceName"      env:"ServiceName"`
	Logger            *logger.LogConfig              `mapstructure:"logger"                              envPrefix:"Logger_"`
	GRPC              *grpc.GrpcConfig               `mapstructure:"grpc"                                envPrefix:"GRPC_"`
	Http              *customEcho.EchoHttpConfig     `mapstructure:"http"                                envPrefix:"Http_"`
	Context           Context                        `mapstructure:"context"                             envPrefix:"Context_"`
	Redis             *redis.RedisConfig             `mapstructure:"redis"                               envPrefix:"Redis_"`
	RabbitMQ          *rabbitmqconfig.RabbitMQConfig `mapstructure:"rabbitmq"                            envPrefix:"RabbitMQ_"`
	OTel              *otel.OpenTelemetryConfig      `mapstructure:"otel"                                envPrefix:"OTel_"`
	OTelMetricsConfig *metrics.OTelMetricsConfig     `mapstructure:"otelMetrics"                         envPrefix:"OTelMetrics_"`
	EventStoreConfig  eventstroredb.EventStoreConfig `mapstructure:"eventStoreConfig"                    envPrefix:"EventStoreConfig_"`
	Elastic           elasticsearch.Config           `mapstructure:"elastic"                             envPrefix:"Elastic_"`
	ElasticIndexes    ElasticIndexes                 `mapstructure:"elasticIndexes"                      envPrefix:"ElasticIndexes_"`
	Mongo             *mongodb.MongoDbConfig         `mapstructure:"mongo"                               envPrefix:"Mongo_"`
	MongoCollections  MongoCollections               `mapstructure:"mongoCollections"                    envPrefix:"MongoCollections_"`
}

type Context struct {
	Timeout int `mapstructure:"timeout" env:"Timeout"`
}

type MongoCollections struct {
	Products string `mapstructure:"products" validate:"required" env:"Products"`
}

type ElasticIndexes struct {
	Products string `mapstructure:"products" validate:"required" env:"Products"`
}

func NewAppConfig(env config.Environment) (*AppConfig, error) {
	cfg, err := config.BindConfig[*AppConfig](env.GetEnvironmentName())
	if err != nil {
		return nil, err
	}

	grpcPort := os.Getenv(constants.GrpcPort)
	if grpcPort != "" {
		cfg.GRPC.Port = grpcPort
	}

	jaegerPort := os.Getenv(constants.JaegerPort)
	if jaegerPort != "" {
		cfg.OTel.JaegerExporterConfig.AgentPort = jaegerPort
	}

	jaegerHost := os.Getenv(constants.JaegerHost)
	if jaegerHost != "" {
		cfg.OTel.JaegerExporterConfig.AgentHost = jaegerHost
	}

	return cfg, nil
}

func (cfg *AppConfig) GetMicroserviceNameUpper() string {
	return strings.ToUpper(cfg.ServiceName)
}

func (cfg *AppConfig) GetMicroserviceName() string {
	return cfg.ServiceName
}
