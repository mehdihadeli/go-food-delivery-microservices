package infrastructure

import (
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel/metric"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/serializer"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/es/contracts"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/eventstroredb"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/grpc"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orders/config"
)

type infrastructureConfigurations struct {
	Log                  logger.Logger
	Cfg                  *config.Config
	Validator            *validator.Validate
	Metrics              metric.Meter
	Esdb                 *esdb.Client
	EsdbSerializer       *eventstroredb.EsdbSerializer
	EventSerializer      serializer.EventSerializer
	CheckpointRepository contracts.SubscriptionCheckpointRepository
	ElasticClient        *elasticsearch.Client
	MongoClient          *mongo.Client
	GrpcClient           grpc.GrpcClient
}
