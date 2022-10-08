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
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel/metric"
	"gorm.io/gorm"
)

type infrastructureConfigurations struct {
	log             logger.Logger
	cfg             *config.Config
	validator       *validator.Validate
	pgConn          *pgxpool.Pool
	gorm            *gorm.DB
	metrics         metric.Meter
	esdb            *esdb.Client
	mongoClient     *mongo.Client
	grpcClient      grpc.GrpcClient
	elasticClient   *elasticsearch.Client
	redis           redis.UniversalClient
	eventSerializer serializer.EventSerializer
}

func (i *infrastructureConfigurations) Log() logger.Logger {
	return i.log
}

func (i *infrastructureConfigurations) Cfg() *config.Config {
	return i.cfg
}

func (i *infrastructureConfigurations) Validator() *validator.Validate {
	return i.validator
}

func (i *infrastructureConfigurations) PgConn() *pgxpool.Pool {
	return i.pgConn
}

func (i *infrastructureConfigurations) Gorm() *gorm.DB {
	return i.gorm
}

func (i *infrastructureConfigurations) Esdb() *esdb.Client {
	return i.esdb
}

func (i *infrastructureConfigurations) MongoClient() *mongo.Client {
	return i.mongoClient
}

func (i *infrastructureConfigurations) GrpcClient() grpc.GrpcClient {
	return i.grpcClient
}

func (i *infrastructureConfigurations) ElasticClient() *elasticsearch.Client {
	return i.elasticClient
}

func (i *infrastructureConfigurations) Redis() redis.UniversalClient {
	return i.redis
}

func (i *infrastructureConfigurations) EventSerializer() serializer.EventSerializer {
	return i.eventSerializer
}

func (i *infrastructureConfigurations) Metrics() metric.Meter {
	return i.metrics
}
