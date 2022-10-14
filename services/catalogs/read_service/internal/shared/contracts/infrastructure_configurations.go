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
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel/metric"
	"gorm.io/gorm"
)

type InfrastructureConfigurations struct {
	Log             logger.Logger
	Cfg             *config.Config
	Validator       *validator.Validate
	PgConn          *pgxpool.Pool
	Gorm            *gorm.DB
	Metrics         metric.Meter
	Esdb            *esdb.Client
	MongoClient     *mongo.Client
	GrpcClient      grpc.GrpcClient
	ElasticClient   *elasticsearch.Client
	Redis           redis.UniversalClient
	EventSerializer serializer.EventSerializer
}
