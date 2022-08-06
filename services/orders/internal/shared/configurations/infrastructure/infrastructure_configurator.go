package infrastructure

import (
	"context"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/go-playground/validator"
	grpcServer "github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	customEcho "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/custom_echo"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/interceptors"
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
	Log               logger.Logger
	Cfg               *config.Config
	Validator         *validator.Validate
	KafkaConn         *kafka.Conn
	KafkaProducer     kafkaClient.Producer
	Im                interceptors.InterceptorManager
	Pgx               *postgres.Pgx
	Gorm              *gorm.DB
	Metrics           *OrdersServiceMetrics
	EchoServer        customEcho.EchoHttpServer
	GrpcServer        grpcServer.GrpcServer
	Esdb              *esdb.Client
	ElasticClient     *v7.Client
	CustomMiddlewares cutomMiddlewares.CustomMiddlewares
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

func NewInfrastructureConfigurator(log logger.Logger, cfg *config.Config, echoServer customEcho.EchoHttpServer, grpcServer grpcServer.GrpcServer) *infrastructureConfigurator {
	return &infrastructureConfigurator{log: log, cfg: cfg, echoServer: echoServer, grpcServer: grpcServer}
}

func (ic *infrastructureConfigurator) ConfigInfrastructures(ctx context.Context) (*InfrastructureConfiguration, error, func()) {

	infrastructure := &InfrastructureConfiguration{Cfg: ic.cfg, EchoServer: ic.echoServer, GrpcServer: ic.grpcServer, Log: ic.log, Validator: validator.New()}

	infrastructure.Im = interceptors.NewInterceptorManager(ic.log)

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

	es, err, eventStoreCleanup := ic.configEventStore()
	if err != nil {
		return nil, err, nil
	}
	cleanup = append(cleanup, eventStoreCleanup)
	infrastructure.Esdb = es

	kafkaConn, kafkaProducer, err, kafkaCleanup := ic.configKafka(ctx)
	if err != nil {
		return nil, err, nil
	}
	cleanup = append(cleanup, kafkaCleanup)
	infrastructure.KafkaConn = kafkaConn
	infrastructure.KafkaProducer = kafkaProducer

	ic.configSwagger()
	ic.configMiddlewares(metrics)
	ic.configureHealthCheckEndpoints(ctx)

	if err != nil {
		return nil, err, nil
	}

	return infrastructure, nil, func() {
		for _, c := range cleanup {
			c()
		}
	}
}
