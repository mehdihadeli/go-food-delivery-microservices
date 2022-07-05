package configurations

import (
	"context"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/go-playground/validator"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/interceptors"
	kafkaClient "github.com/mehdihadeli/store-golang-microservice-sample/pkg/kafka"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/web/middlewares"
	v7 "github.com/olivere/elastic/v7"
	"github.com/opentracing/opentracing-go"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type Infrastructure struct {
	Log               logger.Logger
	Cfg               *config.Config
	Validator         *validator.Validate
	KafkaConn         *kafka.Conn
	KafkaProducer     kafkaClient.Producer
	Im                interceptors.InterceptorManager
	PgConn            *pgxpool.Pool
	Gorm              *gorm.DB
	Metrics           *shared.CatalogsServiceMetrics
	Echo              *echo.Echo
	GrpcServer        *grpc.Server
	Esdb              *esdb.Client
	MongoClient       *mongo.Client
	ElasticClient     *v7.Client
	MiddlewareManager middlewares.MiddlewareManager
}

func (h *Infrastructure) TraceErr(span opentracing.Span, err error) {
	span.SetTag("error", true)
	span.LogKV("error_code", err.Error())
	h.Metrics.ErrorHttpRequests.Inc()
}

var infrastructure *Infrastructure

type InfrastructureConfigurator interface {
	ConfigureInfrastructure() error
}

type infrastructureConfigurator struct {
	log        logger.Logger
	cfg        *config.Config
	echo       *echo.Echo
	grpcServer *grpc.Server
}

func NewInfrastructureConfigurator(log logger.Logger, cfg *config.Config, echo *echo.Echo, grpcServer *grpc.Server) *infrastructureConfigurator {
	return &infrastructureConfigurator{log: log, cfg: cfg, echo: echo, grpcServer: grpcServer}
}

func (ic *infrastructureConfigurator) ConfigInfrastructures(ctx context.Context) (*Infrastructure, error, func()) {

	infrastructure = &Infrastructure{Cfg: ic.cfg, Echo: ic.echo, GrpcServer: ic.grpcServer, Log: ic.log, Validator: validator.New()}

	infrastructure.Im = interceptors.NewInterceptorManager(ic.log)
	infrastructure.Metrics = shared.NewCatalogsServiceMetrics(ic.cfg)

	//infrastructure.MiddlewareManager = middlewares.NewMiddlewareManager(ic.server.Log, ic.server.Cfg, getHttpMetricsCb(infrastructure.Metrics))

	cleanup := []func(){}

	err, jaegerCleanup := ic.configJaeger()
	if err != nil {
		return nil, err, nil
	}
	cleanup = append(cleanup, jaegerCleanup)

	mongoClient, err, mongoCleanup := ic.configMongo(ctx)
	if err != nil {
		return nil, err, nil
	}
	cleanup = append(cleanup, mongoCleanup)
	infrastructure.MongoClient = mongoClient

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
	ic.configMiddlewares()
	ic.configureHealthCheckEndpoints(ctx, mongoClient)

	if err != nil {
		return nil, err, nil
	}

	return infrastructure, nil, func() {
		for _, c := range cleanup {
			defer c()
		}
	}
}
