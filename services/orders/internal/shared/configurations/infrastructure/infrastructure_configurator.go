package infrastructure

import (
	"context"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/go-playground/validator"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/eventstroredb"
	grpcServer "github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	customEcho "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/custom_echo"
	kafkaClient "github.com/mehdihadeli/store-golang-microservice-sample/pkg/kafka"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	postgres "github.com/mehdihadeli/store-golang-microservice-sample/pkg/postgres_pgx"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/web/custom_middlewares"
	v7 "github.com/olivere/elastic/v7"
	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

type InfrastructureConfiguration struct {
	Log                  logger.Logger
	Cfg                  *config.Config
	Validator            *validator.Validate
	KafkaConn            *kafka.Conn
	KafkaProducer        kafkaClient.Producer
	Pgx                  *postgres.Pgx
	Gorm                 *gorm.DB
	Metrics              *OrdersServiceMetrics
	Esdb                 *esdb.Client
	EsdbSerializer       *eventstroredb.EsdbSerializer
	CheckpointRepository contracts.SubscriptionCheckpointRepository
	ElasticClient        *v7.Client
	CustomMiddlewares    cutomMiddlewares.CustomMiddlewares
	EchoServer           customEcho.EchoHttpServer
	GrpcServer           grpcServer.GrpcServer
}

type InfrastructureConfigurator interface {
	ConfigureInfrastructure() error
}

type infrastructureConfigurator struct {
	log        logger.Logger
	cfg        *config.Config
	echoServer customEcho.EchoHttpServer
	grpcServer grpcServer.GrpcServer
}

func NewInfrastructureConfigurator(log logger.Logger, cfg *config.Config, echo customEcho.EchoHttpServer, grpcServer grpcServer.GrpcServer) *infrastructureConfigurator {
	return &infrastructureConfigurator{log: log, cfg: cfg, echoServer: echo, grpcServer: grpcServer}
}

func (ic *infrastructureConfigurator) ConfigInfrastructures(ctx context.Context) (*InfrastructureConfiguration, error, func()) {
	infrastructure := &InfrastructureConfiguration{Cfg: ic.cfg, Log: ic.log, Validator: validator.New(), EchoServer: ic.echoServer, GrpcServer: ic.grpcServer}

	metrics := ic.configCatalogsMetrics()
	infrastructure.Metrics = metrics

	cleanup := []func(){}

	err, jaegerCleanup := ic.configJaeger()
	if err != nil {
		return nil, err, nil
	}
	cleanup = append(cleanup, jaegerCleanup)

	//el, err, _ := ic.configElasticSearch(ctx)
	//if err != nil {
	//	return nil, err, nil
	//}
	//infrastructure.ElasticClient = el

	es, checkpointRepository, esdbSerializer, err, eventStoreCleanup := ic.configEventStore()
	if err != nil {
		return nil, err, nil
	}
	cleanup = append(cleanup, eventStoreCleanup)
	infrastructure.Esdb = es
	infrastructure.CheckpointRepository = checkpointRepository
	infrastructure.EsdbSerializer = esdbSerializer

	kafkaConn, kafkaProducer, err, kafkaCleanup := ic.configKafka(ctx)
	if err != nil {
		return nil, err, nil
	}
	cleanup = append(cleanup, kafkaCleanup)
	infrastructure.KafkaConn = kafkaConn
	infrastructure.KafkaProducer = kafkaProducer

	ic.configSwagger()
	ic.configMiddlewares(infrastructure.Metrics)

	if err != nil {
		return nil, err, nil
	}

	return infrastructure, nil, func() {
		for _, c := range cleanup {
			c()
		}
	}
}
