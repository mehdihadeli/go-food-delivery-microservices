package infrastructure

import (
	"context"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/go-playground/validator"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	grpcServer "github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	customEcho "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/custom_echo"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/interceptors"
	kafkaClient "github.com/mehdihadeli/store-golang-microservice-sample/pkg/kafka"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/web/middlewares"
	v7 "github.com/olivere/elastic/v7"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type InfrastructureConfigurations struct {
	Log               logger.Logger
	Cfg               *config.Config
	Validator         *validator.Validate
	KafkaConn         *kafka.Conn
	KafkaProducer     kafkaClient.Producer
	Im                interceptors.InterceptorManager
	PgConn            *pgxpool.Pool
	Gorm              *gorm.DB
	Metrics           *CatalogsServiceMetrics
	EchoServer        customEcho.EchoHttpServer
	GrpcServer        grpcServer.GrpcServer
	Esdb              *esdb.Client
	MongoClient       *mongo.Client
	ElasticClient     *v7.Client
	Redis             redis.UniversalClient
	MiddlewareManager cutomMiddlewares.CustomMiddlewares
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

func (ic *infrastructureConfigurator) ConfigInfrastructures(ctx context.Context) (*InfrastructureConfigurations, error, func()) {

	infrastructure := &InfrastructureConfigurations{Cfg: ic.cfg, EchoServer: ic.echoServer, GrpcServer: ic.grpcServer, Log: ic.log, Validator: validator.New()}

	infrastructure.Im = interceptors.NewInterceptorManager(ic.log)

	metrics := ic.configCatalogsMetrics()
	infrastructure.Metrics = metrics

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

	redis, err, redisCleanup := ic.configRedis(ctx)
	if err != nil {
		return nil, err, nil
	}
	cleanup = append(cleanup, redisCleanup)
	infrastructure.Redis = redis

	//el, err, _ := ic.configElasticSearch(ctx)
	//if err != nil {
	//	return nil, err, nil
	//}
	//infrastructure.ElasticClient = el

	kafkaConn, kafkaProducer, err, kafkaCleanup := ic.configKafka(ctx)
	if err != nil {
		return nil, err, nil
	}
	cleanup = append(cleanup, kafkaCleanup)
	infrastructure.KafkaConn = kafkaConn
	infrastructure.KafkaProducer = kafkaProducer

	ic.configSwagger()
	ic.configMiddlewares(metrics)
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
