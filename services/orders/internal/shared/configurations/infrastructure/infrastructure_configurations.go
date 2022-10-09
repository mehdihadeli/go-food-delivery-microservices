package infrastructure

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

type infrastructureConfigurations struct {
	log                  logger.Logger
	cfg                  *config.Config
	validator            *validator.Validate
	metrics              metric.Meter
	esdb                 *esdb.Client
	esdbSerializer       *eventstroredb.EsdbSerializer
	eventSerializer      serializer.EventSerializer
	checkpointRepository contracts.SubscriptionCheckpointRepository
	elasticClient        *elasticsearch.Client
	mongoClient          *mongo.Client
	grpcClient           grpc.GrpcClient
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

func (i *infrastructureConfigurations) Esdb() *esdb.Client {
	return i.esdb
}

func (i *infrastructureConfigurations) GrpcClient() grpc.GrpcClient {
	return i.grpcClient
}

func (i *infrastructureConfigurations) ElasticClient() *elasticsearch.Client {
	return i.elasticClient
}

func (i *infrastructureConfigurations) MongoClient() *mongo.Client {
	return i.mongoClient
}

func (i *infrastructureConfigurations) EventSerializer() serializer.EventSerializer {
	return i.eventSerializer
}

func (i *infrastructureConfigurations) EsdbSerializer() *eventstroredb.EsdbSerializer {
	return i.esdbSerializer
}

func (i *infrastructureConfigurations) CheckpointRepository() contracts.SubscriptionCheckpointRepository {
	return i.checkpointRepository
}

func (i *infrastructureConfigurations) Metrics() metric.Meter {
	return i.metrics
}
