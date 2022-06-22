package configurations

import (
	"context"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/heptiolabs/healthcheck"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/constants"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/elasticsearch"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/store"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/eventstroredb"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/interceptors"
	kafkaClient "github.com/mehdihadeli/store-golang-microservice-sample/pkg/kafka"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mongodb"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/postgres"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/consts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/shared"
	catalog_constants "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/shared/constants"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/shared/server"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/shared/web/middlewares"
	v7 "github.com/olivere/elastic/v7"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Infrastructure struct {
	Log               logger.Logger
	Cfg               *config.Config
	HttpServer        *http.Server
	Validator         *validator.Validate
	KafkaConn         *kafka.Conn
	KafkaProducer     kafkaClient.Producer
	Im                interceptors.InterceptorManager
	PgConn            *pgxpool.Pool
	Metrics           *shared.CatalogsServiceMetrics
	Echo              *echo.Echo
	MongoClient       *mongo.Client
	ElasticClient     *v7.Client
	Ctx               context.Context
	MiddlewareManager middlewares.MiddlewareManager
	healthCheck       healthcheck.Handler
}

var infrastructure *Infrastructure

type InfrastructureConfigurator interface {
	ConfigureInfrastructure() error
}

type infrastructureConfigurator struct {
	server *server.Server
}

func NewInfrastructureConfigurator(server *server.Server) *infrastructureConfigurator {
	return &infrastructureConfigurator{server: server}
}

func (ic *infrastructureConfigurator) ConfigInfrastructures(ctx context.Context, cancelFunc context.CancelFunc) (error, *Infrastructure, func()) {

	infrastructure = &Infrastructure{Cfg: ic.server.Cfg, Echo: ic.server.Echo, Log: ic.server.Log}

	infrastructure.Im = interceptors.NewInterceptorManager(ic.server.Log)
	infrastructure.Metrics = shared.NewCatalogsServiceMetrics(ic.server.Cfg)

	infrastructure.MiddlewareManager = middlewares.NewMiddlewareManager(ic.server.Log, ic.server.Cfg, getHttpMetricsCb(infrastructure.Metrics))

	cleanup := []func(){}

	var err, jaegerCleanup = infrastructure.configJaeger()
	if err != nil {
		return err, nil, nil
	}
	cleanup = append(cleanup, jaegerCleanup)

	err, mongoDeferCallback := infrastructure.configMongo(ctx)
	if err != nil {
		return err, nil, nil
	}
	cleanup = append(cleanup, mongoDeferCallback)

	err, postgresCleanup := infrastructure.configPostgres()
	if err != nil {
		return err, nil, nil
	}
	cleanup = append(cleanup, postgresCleanup)

	err, _ = infrastructure.configElasticSearch(ctx)
	if err != nil {
		return err, nil, nil
	}

	err, eventStoreCleanup := infrastructure.configEventStore()
	if err != nil {
		return err, nil, nil
	}
	cleanup = append(cleanup, eventStoreCleanup)

	err, kafkaCleanup := infrastructure.configKafka(ctx)
	if err != nil {
		return err, nil, nil
	}
	cleanup = append(cleanup, kafkaCleanup)

	err, kafkaConsumerCleanup := infrastructure.configKafkaConsumers(ctx)
	if err != nil {
		return err, nil, nil
	}
	cleanup = append(cleanup, kafkaConsumerCleanup)

	infrastructure.configMiddlewares()
	infrastructure.configureHealthCheckEndpoints(ctx)

	if err != nil {
		return err, nil, nil
	}

	return nil, infrastructure, func() {
		for _, deferFunc := range cleanup {
			defer deferFunc()
		}
	}
}

func (i *Infrastructure) configureHealthCheckEndpoints(ctx context.Context) {

	i.healthCheck.AddReadinessCheck(constants.MongoDB, healthcheck.AsyncWithContext(ctx, func() error {
		if err := i.MongoClient.Ping(ctx, nil); err != nil {
			i.Log.Warnf("(MongoDB Readiness Check) err: {%v}", err)
			return err
		}
		return nil
	}, time.Duration(i.Cfg.Probes.CheckIntervalSeconds)*time.Second))

	//i.healthCheck.AddReadinessCheck(constants.ElasticSearch, healthcheck.AsyncWithContext(ctx, func() error {
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

func (i *Infrastructure) configMiddlewares() {

	i.Echo.Use(i.MiddlewareManager.RequestLoggerMiddleware)
	i.Echo.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         catalog_constants.StackSize,
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))
	i.Echo.Use(middleware.RequestID())
	i.Echo.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: catalog_constants.GzipLevel,
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "swagger")
		},
	}))

	i.Echo.Use(middleware.BodyLimit(catalog_constants.BodyLimit))
}

func (i *Infrastructure) configEventStore() (error, func()) {
	db, err := eventstroredb.NewEventStoreDB(i.Cfg.EventStoreConfig)
	if err != nil {
		return err, nil
	}

	aggregateStore := store.NewAggregateStore(i.Log, db)
	fmt.Print(aggregateStore)

	return nil, func() {
		db.Close() // nolint: errcheck
	}
}

func (i *Infrastructure) configKafkaConsumers(ctx context.Context) (error, func()) {
	//productMessageProcessor := NewProductMessageProcessor(s.Log, s.cfg, s.v, s.ps, s.metrics)
	//s.Log.Info("Starting Writer Kafka consumers")
	//cg := kafkaClient.NewConsumerGroup(s.cfg.Kafka.Brokers, s.cfg.Kafka.GroupID, s.Log)
	//go cg.ConsumeTopic(ctx, s.getConsumerGroupTopics(), kafkaConsumer.PoolSize, productMessageProcessor.ProcessMessages)
	//
	return nil, func() {}
}

func (i *Infrastructure) configKafka(ctx context.Context) (error, func()) {

	if err := i.connectKafkaBrokers(ctx); err != nil {
		return errors.Wrap(err, "i.connectKafkaBrokers"), nil
	}

	if i.Cfg.Kafka.InitTopics {
		i.initKafkaTopics(ctx)
	}

	kafkaProducer := kafkaClient.NewProducer(i.Log, i.Cfg.Kafka.Brokers)

	i.KafkaProducer = kafkaProducer

	return nil, func() {
		kafkaProducer.Close() // nolint: errcheck
		i.KafkaConn.Close()   // nolint: errcheck
	}
}

func (i *Infrastructure) configJaeger() (error, func()) {
	if i.Cfg.Jaeger.Enable {
		tracer, closer, err := tracing.NewJaegerTracer(i.Cfg.Jaeger)
		if err != nil {
			return err, nil
		}
		//defer closer.Close() // nolint: errcheck
		opentracing.SetGlobalTracer(tracer)
		return nil, func() { closer.Close() }
	}

	return nil, func() {}
}

func (i *Infrastructure) configMongo(ctx context.Context) (error, func()) {
	mongoDBConn, err := mongodb.NewMongoDBConn(ctx, i.Cfg.Mongo)
	if err != nil {
		return errors.Wrap(err, "NewMongoDBConn"), nil
	}
	i.MongoClient = mongoDBConn
	i.Log.Infof("(Mongo connected) SessionsInProgress: {%v}", mongoDBConn.NumberSessionsInProgress())

	i.initMongoDBCollections(ctx)

	return nil, func() {
		mongoDBConn.Disconnect(ctx) // nolint: errcheck
	}
}

func (i *Infrastructure) configPostgres() (error, func()) {
	pgxConn, err := postgres.NewPgxConn(i.Cfg.Postgresql)
	if err != nil {
		return errors.Wrap(err, "postgresql.NewPgxConn"), nil
	}
	i.PgConn = pgxConn
	i.Log.Infof("postgres connected: %v", pgxConn.Stat().TotalConns())

	return nil, func() {
		pgxConn.Close()
	}
}

func (i *Infrastructure) configElasticSearch(ctx context.Context) (error, func()) {

	elasticClient, err := elasticsearch.NewElasticClient(i.Cfg.Elastic)
	if err != nil {
		return err, nil
	}
	i.ElasticClient = elasticClient

	info, code, err := elasticClient.Ping(i.Cfg.Elastic.URL).Do(ctx)
	if err != nil {
		return errors.Wrap(err, "client.Ping"), nil
	}
	i.Log.Infof("Elasticsearch returned with code {%d} and version {%s}", code, info.Version.Number)

	esVersion, err := elasticClient.ElasticsearchVersion(i.Cfg.Elastic.URL)
	if err != nil {
		return errors.Wrap(err, "client.ElasticsearchVersion"), nil
	}
	i.Log.Infof("Elasticsearch version {%s}", esVersion)

	return nil, func() {}
}

func (i *Infrastructure) connectKafkaBrokers(ctx context.Context) error {
	kafkaConn, err := kafkaClient.NewKafkaConn(ctx, i.Cfg.Kafka)
	if err != nil {
		return errors.Wrap(err, "kafka.NewKafkaCon")
	}

	i.KafkaConn = kafkaConn

	brokers, err := kafkaConn.Brokers()
	if err != nil {
		return errors.Wrap(err, "kafkaConn.Brokers")
	}

	i.Log.Infof("kafka connected to brokers: %+v", brokers)

	return nil
}

func (i *Infrastructure) initKafkaTopics(ctx context.Context) {
	controller, err := i.KafkaConn.Controller()
	if err != nil {
		i.Log.WarnMsg("kafkaConn.Controller", err)
		return
	}

	controllerURI := net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port))
	i.Log.Infof("kafka controller uri: %s", controllerURI)

	conn, err := kafka.DialContext(ctx, "tcp", controllerURI)
	if err != nil {
		i.Log.WarnMsg("initKafkaTopics.DialContext", err)
		return
	}
	defer conn.Close() // nolint: errcheck

	i.Log.Infof("established new kafka controller connection: %s", controllerURI)

	productCreateTopic := kafka.TopicConfig{
		Topic:             i.Cfg.KafkaTopics.ProductCreate.TopicName,
		NumPartitions:     i.Cfg.KafkaTopics.ProductCreate.Partitions,
		ReplicationFactor: i.Cfg.KafkaTopics.ProductCreate.ReplicationFactor,
	}

	productCreatedTopic := kafka.TopicConfig{
		Topic:             i.Cfg.KafkaTopics.ProductCreated.TopicName,
		NumPartitions:     i.Cfg.KafkaTopics.ProductCreated.Partitions,
		ReplicationFactor: i.Cfg.KafkaTopics.ProductCreated.ReplicationFactor,
	}

	productUpdateTopic := kafka.TopicConfig{
		Topic:             i.Cfg.KafkaTopics.ProductUpdate.TopicName,
		NumPartitions:     i.Cfg.KafkaTopics.ProductUpdate.Partitions,
		ReplicationFactor: i.Cfg.KafkaTopics.ProductUpdate.ReplicationFactor,
	}

	productUpdatedTopic := kafka.TopicConfig{
		Topic:             i.Cfg.KafkaTopics.ProductUpdated.TopicName,
		NumPartitions:     i.Cfg.KafkaTopics.ProductUpdated.Partitions,
		ReplicationFactor: i.Cfg.KafkaTopics.ProductUpdated.ReplicationFactor,
	}

	productDeleteTopic := kafka.TopicConfig{
		Topic:             i.Cfg.KafkaTopics.ProductDelete.TopicName,
		NumPartitions:     i.Cfg.KafkaTopics.ProductDelete.Partitions,
		ReplicationFactor: i.Cfg.KafkaTopics.ProductDelete.ReplicationFactor,
	}

	productDeletedTopic := kafka.TopicConfig{
		Topic:             i.Cfg.KafkaTopics.ProductDeleted.TopicName,
		NumPartitions:     i.Cfg.KafkaTopics.ProductDeleted.Partitions,
		ReplicationFactor: i.Cfg.KafkaTopics.ProductDeleted.ReplicationFactor,
	}

	if err := conn.CreateTopics(
		productCreateTopic,
		productUpdateTopic,
		productCreatedTopic,
		productUpdatedTopic,
		productDeleteTopic,
		productDeletedTopic,
	); err != nil {
		i.Log.WarnMsg("kafkaConn.CreateTopics", err)
		return
	}

	i.Log.Infof("kafka topics created or already exists: %+v", []kafka.TopicConfig{productCreateTopic, productUpdateTopic, productCreatedTopic, productUpdatedTopic, productDeleteTopic, productDeletedTopic})
}

func (i *Infrastructure) initMongoDBCollections(ctx context.Context) {
	err := i.MongoClient.Database(i.Cfg.Mongo.Db).CreateCollection(ctx, i.Cfg.MongoCollections.Products)
	if err != nil {
		if !utils.CheckErrMessages(err, catalog_constants.ErrMsgMongoCollectionAlreadyExists) {
			i.Log.Warnf("(CreateCollection) err: {%v}", err)
		}
	}

	indexOptions := options.Index().SetSparse(true).SetUnique(true)
	index, err := i.MongoClient.Database(i.Cfg.Mongo.Db).Collection(i.Cfg.MongoCollections.Products).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: consts.ProductIdIndex, Value: 1}},
		Options: indexOptions,
	})
	if err != nil && !utils.CheckErrMessages(err, catalog_constants.ErrMsgAlreadyExists) {
		i.Log.Warnf("(CreateOne) err: {%v}", err)
	}
	i.Log.Infof("(CreatedIndex) index: {%s}", index)

	list, err := i.MongoClient.Database(i.Cfg.Mongo.Db).Collection(i.Cfg.MongoCollections.Products).Indexes().List(ctx)
	if err != nil {
		i.Log.Warnf("(initMongoDBCollections) [List] err: {%v}", err)
	}

	if list != nil {
		var results []bson.M
		if err := list.All(ctx, &results); err != nil {
			i.Log.Warnf("(All) err: {%v}", err)
		}
		i.Log.Infof("(indexes) results: {%#v}", results)
	}

	collections, err := i.MongoClient.Database(i.Cfg.Mongo.Db).ListCollectionNames(ctx, bson.M{})
	if err != nil {
		i.Log.Warnf("(ListCollections) err: {%v}", err)
	}
	i.Log.Infof("(Collections) created collections: {%v}", collections)
}

func (i *Infrastructure) getConsumerGroupTopics() []string {
	return []string{
		i.Cfg.KafkaTopics.ProductCreate.TopicName,
		i.Cfg.KafkaTopics.ProductUpdate.TopicName,
		i.Cfg.KafkaTopics.ProductDelete.TopicName,
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
