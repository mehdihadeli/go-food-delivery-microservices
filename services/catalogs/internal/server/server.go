package server

import (
	"context"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/store"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/eventstroredb"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/interceptors"
	kafkaClient "github.com/mehdihadeli/store-golang-microservice-sample/pkg/kafka"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mongodb"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/postgres"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/config"
	v1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/delivery/http/v1"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/infrastructure/repositories"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/shared"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/shared/middlewares"
	v7 "github.com/olivere/elastic/v7"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	log           logger.Logger
	cfg           *config.Config
	hs            *http.Server
	v             *validator.Validate
	kafkaConn     *kafka.Conn
	mediator      *mediatr.Mediator
	im            interceptors.InterceptorManager
	pgConn        *pgxpool.Pool
	metrics       *shared.CatalogsServiceMetrics
	echo          *echo.Echo
	mw            middlewares.MiddlewareManager
	mongoClient   *mongo.Client
	elasticClient *v7.Client
	doneCh        chan struct{}
}

func NewServer(log logger.Logger, cfg *config.Config) *Server {
	return &Server{log: log, cfg: cfg, echo: echo.New(), v: validator.New()}
}

func (s *Server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	if s.cfg.Jaeger.Enable {
		tracer, closer, err := tracing.NewJaegerTracer(s.cfg.Jaeger)
		if err != nil {
			return err
		}
		defer closer.Close() // nolint: errcheck
		opentracing.SetGlobalTracer(tracer)
	}

	s.mw = middlewares.NewMiddlewareManager(s.log, s.cfg)
	s.im = interceptors.NewInterceptorManager(s.log)
	s.metrics = shared.NewCatalogsServiceMetrics(s.cfg)

	mongoDBConn, err := mongodb.NewMongoDBConn(ctx, s.cfg.Mongo)
	if err != nil {
		return errors.Wrap(err, "NewMongoDBConn")
	}
	s.mongoClient = mongoDBConn
	defer mongoDBConn.Disconnect(ctx) // nolint: errcheck
	s.log.Infof("(Mongo connected) SessionsInProgress: {%v}", mongoDBConn.NumberSessionsInProgress())

	pgxConn, err := postgres.NewPgxConn(s.cfg.Postgresql)
	if err != nil {
		return errors.Wrap(err, "postgresql.NewPgxConn")
	}
	s.pgConn = pgxConn
	s.log.Infof("postgres connected: %v", pgxConn.Stat().TotalConns())
	defer pgxConn.Close()

	//if err := s.initElasticClient(ctx); err != nil {
	//	s.log.Errorf("(initElasticClient) err: {%v}", err)
	//	return err
	//}

	db, err := eventstroredb.NewEventStoreDB(s.cfg.EventStoreConfig)
	if err != nil {
		return err
	}
	defer db.Close() // nolint: errcheck

	aggregateStore := store.NewAggregateStore(s.log, db)
	fmt.Print(aggregateStore)

	kafkaProducer := kafkaClient.NewProducer(s.log, s.cfg.Kafka.Brokers)
	defer kafkaProducer.Close() // nolint: errcheck

	productRepo := repositories.NewPostgresProductRepository(s.log, s.cfg, pgxConn)

	m, err := shared.NewMediator(s.log, s.cfg, productRepo, kafkaProducer)

	if err != nil {
		return err
	}

	s.mediator = m

	//productMessageProcessor := kafkaConsumer.NewProductMessageProcessor(s.log, s.cfg, s.v, s.ps, s.metrics)
	//s.log.Info("Starting Writer Kafka consumers")
	//cg := kafkaClient.NewConsumerGroup(s.cfg.Kafka.Brokers, s.cfg.Kafka.GroupID, s.log)
	//go cg.ConsumeTopic(ctx, s.getConsumerGroupTopics(), kafkaConsumer.PoolSize, productMessageProcessor.ProcessMessages)

	s.initMongoDBCollections(ctx)
	s.runMetrics(cancel)
	s.runHealthCheck(ctx)

	//s.applyVersioningFromHeader()
	productHandlers := v1.NewProductsHandlers(s.echo, s.log, s.mw, s.cfg, s.mediator, s.v, s.metrics)
	productHandlers.MapRoutes()

	go func() {
		if err := s.runHttpServer(); err != nil {
			s.log.Errorf("(s.runHttpServer) err: {%v}", err)
			cancel()
		}
	}()
	s.log.Infof("%s is listening on PORT: {%s}", GetMicroserviceName(s.cfg), s.cfg.Http.Port)

	closeGrpcServer, grpcServer, err := s.newCatalogsServiceGrpcServer()
	if err != nil {
		return errors.Wrap(err, "NewScmGrpcServer")
	}
	defer closeGrpcServer() // nolint: errcheck

	if err := s.connectKafkaBrokers(ctx); err != nil {
		return errors.Wrap(err, "s.connectKafkaBrokers")
	}
	defer s.kafkaConn.Close() // nolint: errcheck

	if s.cfg.Kafka.InitTopics {
		s.initKafkaTopics(ctx)
	}

	<-ctx.Done()
	s.waitShootDown(waitShotDownDuration)
	grpcServer.GracefulStop()

	if err := s.shutDownHealthCheckServer(ctx); err != nil {
		s.log.Warnf("(shutDownHealthCheckServer) err: {%v}", err)
	}
	if err := s.echo.Shutdown(ctx); err != nil {
		s.log.Warnf("(Shutdown) err: {%v}", err)
	}

	<-s.doneCh
	s.log.Infof("%s server exited properly", GetMicroserviceName(s.cfg))

	return nil
}
