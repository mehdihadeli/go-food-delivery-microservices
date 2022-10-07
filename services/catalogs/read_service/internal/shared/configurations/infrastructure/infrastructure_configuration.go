package infrastructure

import (
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-playground/validator"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/serializer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	messageBus "github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/bus"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/producer"
	rabbitmqConfigurations "github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/configurations"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/config"
	cutomMiddlewares "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/web/middlewares"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type infrastructureConfigurations struct {
	Log                          logger.Logger
	Cfg                          *config.Config
	Validator                    *validator.Validate
	Producer                     producer.Producer
	PgConn                       *pgxpool.Pool
	Gorm                         *gorm.DB
	Metrics                      *CatalogsServiceMetrics
	Esdb                         *esdb.Client
	MongoClient                  *mongo.Client
	GrpcClient                   grpc.GrpcClient
	ElasticClient                *elasticsearch.Client
	Redis                        redis.UniversalClient
	MiddlewareManager            cutomMiddlewares.CustomMiddlewares
	EventSerializer              serializer.EventSerializer
	RabbitMQConfigurationBuilder rabbitmqConfigurations.RabbitMQConfigurationBuilder
	RabbitMQBus                  messageBus.Bus
}

func (i *infrastructureConfigurations) GetLog() logger.Logger {
	return i.Log
}

func (i *infrastructureConfigurations) GetCfg() *config.Config {
	return i.Cfg
}

func (i *infrastructureConfigurations) GetValidator() *validator.Validate {
	return i.Validator
}

func (i *infrastructureConfigurations) GetProducer() producer.Producer {
	return i.Producer
}

func (i *infrastructureConfigurations) GetPgConn() *pgxpool.Pool {
	return i.PgConn
}

func (i *infrastructureConfigurations) GetGorm() *gorm.DB {
	return i.Gorm
}

func (i *infrastructureConfigurations) GetEsdb() *esdb.Client {
	return i.Esdb
}

func (i *infrastructureConfigurations) GetMongoClient() *mongo.Client {
	return i.MongoClient
}

func (i *infrastructureConfigurations) GetGrpcClient() grpc.GrpcClient {
	return i.GrpcClient
}

func (i *infrastructureConfigurations) GetElasticClient() *elasticsearch.Client {
	return i.ElasticClient
}

func (i *infrastructureConfigurations) GetRedis() redis.UniversalClient {
	return i.Redis
}

func (i *infrastructureConfigurations) GetMiddlewareManager() cutomMiddlewares.CustomMiddlewares {
	return i.MiddlewareManager
}

func (i *infrastructureConfigurations) GetEventSerializer() serializer.EventSerializer {
	return i.EventSerializer
}

func (i *infrastructureConfigurations) GetRabbitMQConfigurationBuilder() rabbitmqConfigurations.RabbitMQConfigurationBuilder {
	return i.RabbitMQConfigurationBuilder
}

func (i *infrastructureConfigurations) GetRabbitMQBus() messageBus.Bus {
	return i.RabbitMQBus
}
