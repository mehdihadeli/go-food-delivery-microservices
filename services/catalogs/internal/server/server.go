package server

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/interceptors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mongodb"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/config"
	v7 "github.com/olivere/elastic/v7"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

type server struct {
	cfg           *config.Config
	log           logger.Logger
	im            interceptors.InterceptorManager
	mw            middlewares.MiddlewareManager
	os            *service.OrderService
	v             *validator.Validate
	mongoClient   *mongo.Client
	elasticClient *v7.Client
	echo          *echo.Echo
	metrics       *metrics.ESMicroserviceMetrics
	ps            *http.Server
	doneCh        chan struct{}
}

func NewServer(cfg *config.Config, log logger.Logger) *server {
	return &server{cfg: cfg, log: log, v: validator.New(), echo: echo.New(), doneCh: make(chan struct{})}
}

func (s *server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	if err := s.v.StructCtx(ctx, s.cfg); err != nil {
		return errors.Wrap(err, "cfg validate")
	}

	if s.cfg.Jaeger.Enable {
		tracer, closer, err := tracing.NewJaegerTracer(s.cfg.Jaeger)
		if err != nil {
			return err
		}
		defer closer.Close() // nolint: errcheck
		opentracing.SetGlobalTracer(tracer)
	}

	s.metrics = metrics.NewESMicroserviceMetrics(s.cfg)
	s.im = interceptors.NewInterceptorManager(s.log, s.getGrpcMetricsCb())
	s.mw = middlewares.NewMiddlewareManager(s.log, s.cfg, s.getHttpMetricsCb())

	mongoDBConn, err := mongodb.NewMongoDBConn(ctx, s.cfg.Mongo)
	if err != nil {
		return errors.Wrap(err, "NewMongoDBConn")
	}
	s.mongoClient = mongoDBConn
	defer mongoDBConn.Disconnect(ctx) // nolint: errcheck
	s.log.Infof("(Mongo connected) SessionsInProgress: {%v}", mongoDBConn.NumberSessionsInProgress())

	if err := s.initElasticClient(ctx); err != nil {
		s.log.Errorf("(initElasticClient) err: {%v}", err)
		return err
	}

	mongoRepository := repository.NewMongoRepository(s.log, s.cfg, s.mongoClient)
	elasticRepository := repository.NewElasticRepository(s.log, s.cfg, s.elasticClient)

	db, err := eventstroredb.NewEventStoreDB(s.cfg.EventStoreConfig)
	if err != nil {
		return err
	}
	defer db.Close() // nolint: errcheck

	aggregateStore := store.NewAggregateStore(s.log, db)
	s.os = service.NewOrderService(s.log, s.cfg, aggregateStore, mongoRepository, elasticRepository)

	mongoProjection := mongo_projection.NewOrderProjection(s.log, db, mongoRepository, s.cfg)
	elasticProjection := elastic_projection.NewElasticProjection(s.log, db, elasticRepository, s.cfg)

	go func() {
		err := mongoProjection.Subscribe(ctx, []string{s.cfg.Subscriptions.OrderPrefix}, s.cfg.Subscriptions.PoolSize, mongoProjection.ProcessEvents)
		if err != nil {
			s.log.Errorf("(orderProjection.Subscribe) err: {%v}", err)
			cancel()
		}
	}()

	go func() {
		err := elasticProjection.Subscribe(ctx, []string{s.cfg.Subscriptions.OrderPrefix}, s.cfg.Subscriptions.PoolSize, elasticProjection.ProcessEvents)
		if err != nil {
			s.log.Errorf("(elasticProjection.Subscribe) err: {%v}", err)
			cancel()
		}
	}()

	orderHandlers := orderHttp.NewOrderHandlers(s.echo.Group(s.cfg.Http.OrdersPath), s.log, s.mw, s.cfg, s.v, s.os, s.metrics)
	orderHandlers.MapRoutes()

	s.initMongoDBCollections(ctx)
	s.runMetrics(cancel)
	s.runHealthCheck(ctx)

	go func() {
		if err := s.runHttpServer(); err != nil {
			s.log.Errorf("(s.runHttpServer) err: {%v}", err)
			cancel()
		}
	}()
	s.log.Infof("%s is listening on PORT: {%s}", GetMicroserviceName(s.cfg), s.cfg.Http.Port)

	closeGrpcServer, grpcServer, err := s.newOrderGrpcServer()
	if err != nil {
		cancel()
		return err
	}
	defer closeGrpcServer() // nolint: errcheck

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
