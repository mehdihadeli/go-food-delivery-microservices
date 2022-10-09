package contracts

import (
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-playground/validator"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/serializer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/eventstroredb"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel/metric"
)

type InfrastructureConfigurations interface {
	Log() logger.Logger
	Cfg() *config.Config
	Validator() *validator.Validate
	Esdb() *esdb.Client
	GrpcClient() grpc.GrpcClient
	ElasticClient() *elasticsearch.Client
	MongoClient() *mongo.Client
	EventSerializer() serializer.EventSerializer
	EsdbSerializer() *eventstroredb.EsdbSerializer
	CheckpointRepository() contracts.SubscriptionCheckpointRepository
	Metrics() metric.Meter
}
