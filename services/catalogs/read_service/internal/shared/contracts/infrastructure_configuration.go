package contracts

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

type InfrastructureConfiguration interface {
	GetLog() logger.Logger
	GetCfg() *config.Config
	GetValidator() *validator.Validate
	GetProducer() producer.Producer
	GetPgConn() *pgxpool.Pool
	GetGorm() *gorm.DB
	GetEsdb() *esdb.Client
	GetMongoClient() *mongo.Client
	GetGrpcClient() grpc.GrpcClient
	GetElasticClient() *elasticsearch.Client
	GetRedis() redis.UniversalClient
	GetMiddlewareManager() cutomMiddlewares.CustomMiddlewares
	GetEventSerializer() serializer.EventSerializer
	GetRabbitMQConfigurationBuilder() rabbitmqConfigurations.RabbitMQConfigurationBuilder
	GetRabbitMQBus() messageBus.Bus
}
