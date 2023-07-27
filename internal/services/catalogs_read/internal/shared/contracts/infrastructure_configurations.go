package contracts

import (
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-playground/validator"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/serializer"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/grpc"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel/metric"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/config"
)

type InfrastructureConfigurations struct {
	Log             logger.Logger
	Cfg             *config.Config
	Validator       *validator.Validate
	PgConn          *pgxpool.Pool
	Metrics         metric.Meter
	Esdb            *esdb.Client
	MongoClient     *mongo.Client
	GrpcClient      grpc.GrpcClient
	ElasticClient   *elasticsearch.Client
	Redis           redis.UniversalClient
	EventSerializer serializer.EventSerializer
}
