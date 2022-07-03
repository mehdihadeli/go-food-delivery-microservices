package configurations

import (
	"context"
	"fmt"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/go-playground/validator"
	"github.com/heptiolabs/healthcheck"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/constants"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/elasticsearch"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/store"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/eventstroredb"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/gorm_postgres"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/interceptors"
	kafkaClient "github.com/mehdihadeli/store-golang-microservice-sample/pkg/kafka"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mongodb"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/postgres"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/docs"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/consts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/models"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared"
	catalog_constants "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/constants"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/web/middlewares"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/web/middlewares/problem_details"
	v7 "github.com/olivere/elastic/v7"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"net"
	"strconv"
	"strings"
	"time"
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

	gorm, err := ic.configGorm()
	if err != nil {
		return nil, err, nil
	}
	infrastructure.Gorm = gorm

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

	pgx, err, postgresCleanup := ic.configPostgres()
	if err != nil {
		return nil, err, nil
	}
	cleanup = append(cleanup, postgresCleanup)
	infrastructure.PgConn = pgx

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

	err, kafkaConsumerCleanup := ic.configKafkaConsumers(ctx)
	if err != nil {
		return nil, err, nil
	}
	cleanup = append(cleanup, kafkaConsumerCleanup)

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

func (ic *infrastructureConfigurator) configureHealthCheckEndpoints(ctx context.Context, mongoClient *mongo.Client) {

	health := healthcheck.NewHandler()

	health.AddReadinessCheck(constants.MongoDB, healthcheck.AsyncWithContext(ctx, func() error {
		if err := mongoClient.Ping(ctx, nil); err != nil {
			ic.log.Warnf("(MongoDB Readiness Check) err: {%v}", err)
			return err
		}
		return nil
	}, time.Duration(ic.cfg.Probes.CheckIntervalSeconds)*time.Second))

	//health.AddReadinessCheck(constants.ElasticSearch, healthcheck.AsyncWithContext(ctx, func() error {
	//	_, _, err := s.elasticClient.Ping(s.cfg.Elastic.URL).Do(ctx)
	//	if err != nil {
	//		s.log.Warnf("(ElasticSearch Readiness Check) err: {%v}", err)
	//		return errors.Wrap(err, "client.Ping")
	//	}
	//	return nil
	//}, time.Duration(s.cfg.Probes.CheckIntervalSeconds)*time.Second))
	//
	//i.healthCheck.AddLivenessCheck(constants.ElasticSearch, healthcheck.AsyncWithContext(ctx, func() error {
	//	_, _, err := s.elasticClient.Ping(s.cfg.Elastic.URL).Do(ctx)
	//	if err != nil {
	//		s.log.Warnf("(ElasticSearch Liveness Check) err: {%v}", err)
	//		return errors.Wrap(err, "client.Ping")
	//	}
	//	return nil
	//}, time.Duration(s.cfg.Probes.CheckIntervalSeconds)*time.Second))
}

func (ic *infrastructureConfigurator) configMiddlewares() {

	ic.echo.HideBanner = true

	ic.echo.HTTPErrorHandler = problem_details.ProblemHandler

	//i.Echo.Use(i.MiddlewareManager.RequestLoggerMiddleware)
	ic.echo.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         catalog_constants.StackSize,
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))
	ic.echo.Use(middleware.RequestID())
	ic.echo.Use(middleware.Logger())
	ic.echo.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: catalog_constants.GzipLevel,
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "swagger")
		},
	}))

	ic.echo.Use(middleware.BodyLimit(catalog_constants.BodyLimit))
}

func (ic *infrastructureConfigurator) configGorm() (*gorm.DB, error) {
	gorm, err := gorm_postgres.NewGorm(ic.cfg.GormPostgres)
	if err != nil {
		return nil, err
	}

	err = gorm.AutoMigrate(&models.Product{})
	if err != nil {
		return nil, err
	}

	return gorm, nil
}

func (ic *infrastructureConfigurator) configSwagger() {
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Title = "Catalogs Service Api"
	docs.SwaggerInfo.Description = "Catalogs Service Api."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api/v1"

	ic.echo.GET("/swagger/*", echoSwagger.WrapHandler)
}

func (ic *infrastructureConfigurator) configEventStore() (*esdb.Client, error, func()) {
	db, err := eventstroredb.NewEventStoreDB(ic.cfg.EventStoreConfig)
	if err != nil {
		return nil, err, nil
	}

	aggregateStore := store.NewAggregateStore(ic.log, db)
	fmt.Print(aggregateStore)

	return db, nil, func() {
		_ = db.Close() // nolint: errcheck
	}
}

func (ic *infrastructureConfigurator) configKafkaConsumers(ctx context.Context) (error, func()) {
	//productMessageProcessor := NewProductMessageProcessor(s.Log, s.cfg, s.v, s.ps, s.metrics)
	//s.Log.Info("Starting Writer Kafka consumers")
	//cg := kafkaClient.NewConsumerGroup(s.cfg.Kafka.Brokers, s.cfg.Kafka.GroupID, s.Log)
	//go cg.ConsumeTopic(ctx, s.getConsumerGroupTopics(), kafkaConsumer.PoolSize, productMessageProcessor.ProcessMessages)
	//
	return nil, func() {}
}

func (ic *infrastructureConfigurator) configKafka(ctx context.Context) (*kafka.Conn, kafkaClient.Producer, error, func()) {

	kafkaConn, err, kafkaConnCleanup := ic.connectKafkaBrokers(ctx)

	if err != nil {
		return nil, nil, errors.Wrap(err, "i.connectKafkaBrokers"), nil
	}

	if ic.cfg.Kafka.InitTopics {
		ic.initKafkaTopics(ctx, kafkaConn)
	}

	kafkaProducer := kafkaClient.NewProducer(ic.log, ic.cfg.Kafka.Brokers)

	return kafkaConn, kafkaProducer, nil, func() {
		_ = kafkaProducer.Close() // nolint:
		kafkaConnCleanup()
	}
}

func (ic *infrastructureConfigurator) configJaeger() (error, func()) {
	if ic.cfg.Jaeger.Enable {
		tracer, closer, err := tracing.NewJaegerTracer(ic.cfg.Jaeger)
		if err != nil {
			return err, nil
		}
		opentracing.SetGlobalTracer(tracer)
		return nil, func() {
			_ = closer.Close()
		}
	}

	return nil, func() {}
}

func (ic *infrastructureConfigurator) configMongo(ctx context.Context) (*mongo.Client, error, func()) {
	mongoClient, err := mongodb.NewMongoDBConn(ctx, ic.cfg.Mongo)
	if err != nil {
		return nil, errors.Wrap(err, "NewMongoDBConn"), nil
	}

	ic.log.Infof("(Mongo connected) SessionsInProgress: {%v}", mongoClient.NumberSessionsInProgress())

	ic.initMongoDBCollections(ctx, mongoClient)

	return mongoClient, nil, func() {
		_ = mongoClient.Disconnect(ctx) // nolint: errcheck
	}
}

func (ic *infrastructureConfigurator) configPostgres() (*pgxpool.Pool, error, func()) {
	pgxConn, err := postgres.NewPgxConn(ic.cfg.Postgresql)
	if err != nil {
		return nil, errors.Wrap(err, "postgresql.NewPgxConn"), nil
	}

	ic.log.Infof("postgres connected: %v", pgxConn.Stat().TotalConns())

	return pgxConn, nil, func() {
		pgxConn.Close()
	}
}

func (ic *infrastructureConfigurator) configElasticSearch(ctx context.Context) (*v7.Client, error, func()) {

	elasticClient, err := elasticsearch.NewElasticClient(ic.cfg.Elastic)
	if err != nil {
		return nil, err, nil
	}

	info, code, err := elasticClient.Ping(ic.cfg.Elastic.URL).Do(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "client.Ping"), nil
	}
	ic.log.Infof("Elasticsearch returned with code {%d} and version {%s}", code, info.Version.Number)

	esVersion, err := elasticClient.ElasticsearchVersion(ic.cfg.Elastic.URL)
	if err != nil {
		return nil, errors.Wrap(err, "client.ElasticsearchVersion"), nil
	}
	ic.log.Infof("Elasticsearch version {%s}", esVersion)

	return elasticClient, nil, func() {}
}

func (ic *infrastructureConfigurator) connectKafkaBrokers(ctx context.Context) (*kafka.Conn, error, func()) {
	kafkaConn, err := kafkaClient.NewKafkaConn(ctx, ic.cfg.Kafka)
	if err != nil {
		return nil, errors.Wrap(err, "kafka.NewKafkaCon"), nil
	}

	brokers, err := kafkaConn.Brokers()
	if err != nil {
		return nil, errors.Wrap(err, "kafkaConn.Brokers"), nil
	}

	ic.log.Infof("kafka connected to brokers: %+v", brokers)

	return kafkaConn, nil, func() {
		_ = kafkaConn.Close() // nolint: errcheck
	}
}

func (ic *infrastructureConfigurator) initKafkaTopics(ctx context.Context, kafkaConn *kafka.Conn) {
	controller, err := kafkaConn.Controller()
	if err != nil {
		ic.log.WarnMsg("kafkaConn.Controller", err)
		return
	}

	controllerURI := net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port))
	ic.log.Infof("kafka controller uri: %s", controllerURI)

	conn, err := kafka.DialContext(ctx, "tcp", controllerURI)
	if err != nil {
		ic.log.WarnMsg("initKafkaTopics.DialContext", err)
		return
	}
	defer conn.Close() // nolint: errcheck

	ic.log.Infof("established new kafka controller connection: %s", controllerURI)

	productCreateTopic := kafka.TopicConfig{
		Topic:             ic.cfg.KafkaTopics.ProductCreate.TopicName,
		NumPartitions:     ic.cfg.KafkaTopics.ProductCreate.Partitions,
		ReplicationFactor: ic.cfg.KafkaTopics.ProductCreate.ReplicationFactor,
	}

	productCreatedTopic := kafka.TopicConfig{
		Topic:             ic.cfg.KafkaTopics.ProductCreated.TopicName,
		NumPartitions:     ic.cfg.KafkaTopics.ProductCreated.Partitions,
		ReplicationFactor: ic.cfg.KafkaTopics.ProductCreated.ReplicationFactor,
	}

	productUpdateTopic := kafka.TopicConfig{
		Topic:             ic.cfg.KafkaTopics.ProductUpdate.TopicName,
		NumPartitions:     ic.cfg.KafkaTopics.ProductUpdate.Partitions,
		ReplicationFactor: ic.cfg.KafkaTopics.ProductUpdate.ReplicationFactor,
	}

	productUpdatedTopic := kafka.TopicConfig{
		Topic:             ic.cfg.KafkaTopics.ProductUpdated.TopicName,
		NumPartitions:     ic.cfg.KafkaTopics.ProductUpdated.Partitions,
		ReplicationFactor: ic.cfg.KafkaTopics.ProductUpdated.ReplicationFactor,
	}

	productDeleteTopic := kafka.TopicConfig{
		Topic:             ic.cfg.KafkaTopics.ProductDelete.TopicName,
		NumPartitions:     ic.cfg.KafkaTopics.ProductDelete.Partitions,
		ReplicationFactor: ic.cfg.KafkaTopics.ProductDelete.ReplicationFactor,
	}

	productDeletedTopic := kafka.TopicConfig{
		Topic:             ic.cfg.KafkaTopics.ProductDeleted.TopicName,
		NumPartitions:     ic.cfg.KafkaTopics.ProductDeleted.Partitions,
		ReplicationFactor: ic.cfg.KafkaTopics.ProductDeleted.ReplicationFactor,
	}

	if err := conn.CreateTopics(
		productCreateTopic,
		productUpdateTopic,
		productCreatedTopic,
		productUpdatedTopic,
		productDeleteTopic,
		productDeletedTopic,
	); err != nil {
		ic.log.WarnMsg("kafkaConn.CreateTopics", err)
		return
	}

	ic.log.Infof("kafka topics created or already exists: %+v", []kafka.TopicConfig{productCreateTopic, productUpdateTopic, productCreatedTopic, productUpdatedTopic, productDeleteTopic, productDeletedTopic})
}

func (ic *infrastructureConfigurator) initMongoDBCollections(ctx context.Context, mongoClient *mongo.Client) {
	err := mongoClient.Database(ic.cfg.Mongo.Db).CreateCollection(ctx, ic.cfg.MongoCollections.Products)
	if err != nil {
		if !utils.CheckErrMessages(err, catalog_constants.ErrMsgMongoCollectionAlreadyExists) {
			ic.log.Warnf("(CreateCollection) err: {%v}", err)
		}
	}

	indexOptions := options.Index().SetSparse(true).SetUnique(true)
	index, err := mongoClient.Database(ic.cfg.Mongo.Db).Collection(ic.cfg.MongoCollections.Products).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: consts.ProductIdIndex, Value: 1}},
		Options: indexOptions,
	})
	if err != nil && !utils.CheckErrMessages(err, catalog_constants.ErrMsgAlreadyExists) {
		ic.log.Warnf("(CreateOne) err: {%v}", err)
	}
	ic.log.Infof("(CreatedIndex) index: {%s}", index)

	list, err := mongoClient.Database(ic.cfg.Mongo.Db).Collection(ic.cfg.MongoCollections.Products).Indexes().List(ctx)
	if err != nil {
		ic.log.Warnf("(initMongoDBCollections) [List] err: {%v}", err)
	}

	if list != nil {
		var results []bson.M
		if err := list.All(ctx, &results); err != nil {
			ic.log.Warnf("(All) err: {%v}", err)
		}
		ic.log.Infof("(indexes) results: {%#v}", results)
	}

	collections, err := mongoClient.Database(ic.cfg.Mongo.Db).ListCollectionNames(ctx, bson.M{})
	if err != nil {
		ic.log.Warnf("(ListCollections) err: {%v}", err)
	}
	ic.log.Infof("(Collections) created collections: {%v}", collections)
}

func (ic *infrastructureConfigurator) getConsumerGroupTopics() []string {
	return []string{
		ic.cfg.KafkaTopics.ProductCreate.TopicName,
		ic.cfg.KafkaTopics.ProductUpdate.TopicName,
		ic.cfg.KafkaTopics.ProductDelete.TopicName,
	}
}

func getHttpMetricsCb(metrics *shared.CatalogsServiceMetrics) func(err error) {
	return func(err error) {
		if err != nil {
			metrics.ErrorHttpRequests.Inc()
		} else {
			metrics.SuccessHttpRequests.Inc()
		}
	}
}

func getGrpcMetricsCb(metrics *shared.CatalogsServiceMetrics) func(err error) {
	return func(err error) {
		if err != nil {
			metrics.ErrorGrpcRequests.Inc()
		} else {
			metrics.SuccessGrpcRequests.Inc()
		}
	}
}
