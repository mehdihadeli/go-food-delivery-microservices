package configurations

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/elasticsearch"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/store"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/eventstroredb"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/interceptors"
	kafkaClient "github.com/mehdihadeli/store-golang-microservice-sample/pkg/kafka"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mongodb"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/postgres"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/consts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/shared"
	catalog_constants "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/shared/constants"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/shared/server"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net"
	"strconv"
)

//type InfrastructureConfigurator interface {
//	ConfigInfrastructures(ctx context.Context) (error, *InfrastructureConfiguration, func())
//	GetConfiguration() *InfrastructureConfiguration
//}
//
//type infrastructureConfigurator struct {
//	server *server.Server
//	ctx    context.Context
//}
//
//type InfrastructureConfiguration struct {
//	Log           logger.Logger
//	Cfg           *config.Config
//	HttpServer    *http.Server
//	Validator     *validator.Validate
//	KafkaConn     *kafka.Conn
//	KafkaProducer kafkaClient.Producer
//	Im            interceptors.InterceptorManager
//	PgConn        *pgxpool.Pool
//	Metrics       *shared.CatalogsServiceMetrics
//	Echo          *echo.Echo
//	MongoClient   *mongo.Client
//	ElasticClient *v7.Client
//	Ctx           context.Context
//}
//
//var configuration *InfrastructureConfiguration
//
//func NewInfrastructureConfigurator(s *server.Server, ctx context.Context) *infrastructureConfigurator {
//	return &infrastructureConfigurator{server: s, ctx: ctx}
//}

//func (s *infrastructureConfigurator) GetConfiguration() *InfrastructureConfiguration {
//	return configuration
//}

type InfrastructureConfigurator interface {
	ConfigureInfrastructure() error
}

type infrastructureConfigurator struct {
	server *server.Server
}

func NewInfrastructureConfigurator(server *server.Server) *infrastructureConfigurator {
	return &infrastructureConfigurator{server: server}
}

func (ic *infrastructureConfigurator) ConfigInfrastructures(ctx context.Context, cancelFunc context.CancelFunc) (error, func()) {

	ic.server.Im = interceptors.NewInterceptorManager(ic.server.Log)
	ic.server.Metrics = shared.NewCatalogsServiceMetrics(ic.server.Cfg)

	defers := []func(){}

	defers = append(defers, func() {})
	defers = append(defers, func() {})

	var err, jaegerDeferCallback = s.configJaeger()
	if err != nil {
		return err, nil
	}
	defers = append(defers, jaegerDeferCallback)

	err, mongoDeferCallback := s.configMongo(ctx)
	if err != nil {
		return err, nil
	}
	defers = append(defers, mongoDeferCallback)

	err, postgressDeferCallback := s.configPostgres()
	if err != nil {
		return err, nil
	}
	defers = append(defers, postgressDeferCallback)

	err, _ = s.configElasticSearch(ctx)
	if err != nil {
		return err, nil
	}

	err, eventStoreDeferCallback := s.configEventStore()
	if err != nil {
		return err, nil
	}
	defers = append(defers, eventStoreDeferCallback)

	err, kafkaDeferCallback := s.configKafka(ctx)
	if err != nil {
		return err, nil
	}
	defers = append(defers, kafkaDeferCallback)

	err, kafkaConsumerDeferCallback := s.configKafkaConsumers(ctx)
	if err != nil {
		return err, nil
	}
	defers = append(defers, kafkaConsumerDeferCallback)

	err, grpcDeferCallback := s.newCatalogsServiceGrpcServer()
	if err != nil {
		return err, nil
	}
	defers = append(defers, grpcDeferCallback)

	s.initMongoDBCollections(ctx)

	s.runMetrics(cancelFunc)

	s.runHealthCheck(ctx)

	if err != nil {
		return err, nil
	}

	return nil, func() {
		for _, deferFunc := range defers {
			defer deferFunc()
		}
	}
}

func (ic *infrastructureConfigurator) configEventStore() (error, func()) {
	db, err := eventstroredb.NewEventStoreDB(s.Cfg.EventStoreConfig)
	if err != nil {
		return err, nil
	}

	aggregateStore := store.NewAggregateStore(s.Log, db)
	fmt.Print(aggregateStore)

	return nil, func() {
		db.Close() // nolint: errcheck
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

func (ic *infrastructureConfigurator) configKafka(ctx context.Context) (error, func()) {

	if err := ic.connectKafkaBrokers(ctx); err != nil {
		return errors.Wrap(err, "s.connectKafkaBrokers"), nil
	}

	if s.Cfg.Kafka.InitTopics {
		s.initKafkaTopics(ctx)
	}

	kafkaProducer := kafkaClient.NewProducer(s.Log, s.Cfg.Kafka.Brokers)

	s.KafkaProducer = kafkaProducer

	return nil, func() {
		kafkaProducer.Close() // nolint: errcheck
		s.KafkaConn.Close()   // nolint: errcheck
	}
}

func (ic *infrastructureConfigurator) configJaeger() (error, func()) {
	if s.Cfg.Jaeger.Enable {
		tracer, closer, err := tracing.NewJaegerTracer(s.Cfg.Jaeger)
		if err != nil {
			return err, nil
		}
		//defer closer.Close() // nolint: errcheck
		opentracing.SetGlobalTracer(tracer)
		return nil, func() { closer.Close() }
	}

	return nil, func() {}
}

func (s *Server) configMongo(ctx context.Context) (error, func()) {
	mongoDBConn, err := mongodb.NewMongoDBConn(ctx, s.Cfg.Mongo)
	if err != nil {
		return errors.Wrap(err, "NewMongoDBConn"), nil
	}
	s.MongoClient = mongoDBConn
	s.Log.Infof("(Mongo connected) SessionsInProgress: {%v}", mongoDBConn.NumberSessionsInProgress())

	s.initMongoDBCollections(ctx)

	return nil, func() {
		mongoDBConn.Disconnect(ctx) // nolint: errcheck
	}
}

func (s *Server) configPostgres() (error, func()) {
	pgxConn, err := postgres.NewPgxConn(s.Cfg.Postgresql)
	if err != nil {
		return errors.Wrap(err, "postgresql.NewPgxConn"), nil
	}
	s.PgConn = pgxConn
	s.Log.Infof("postgres connected: %v", pgxConn.Stat().TotalConns())

	return nil, func() {
		pgxConn.Close()
	}
}

func (s *Server) configElasticSearch(ctx context.Context) (error, func()) {

	elasticClient, err := elasticsearch.NewElasticClient(s.Cfg.Elastic)
	if err != nil {
		return err, nil
	}
	s.ElasticClient = elasticClient

	info, code, err := elasticClient.Ping(s.Cfg.Elastic.URL).Do(ctx)
	if err != nil {
		return errors.Wrap(err, "client.Ping"), nil
	}
	s.Log.Infof("Elasticsearch returned with code {%d} and version {%s}", code, info.Version.Number)

	esVersion, err := elasticClient.ElasticsearchVersion(s.Cfg.Elastic.URL)
	if err != nil {
		return errors.Wrap(err, "client.ElasticsearchVersion"), nil
	}
	s.Log.Infof("Elasticsearch version {%s}", esVersion)

	return nil, func() {}
}

func (s *Server) connectKafkaBrokers(ctx context.Context) error {
	kafkaConn, err := kafkaClient.NewKafkaConn(ctx, s.Cfg.Kafka)
	if err != nil {
		return errors.Wrap(err, "kafka.NewKafkaCon")
	}

	s.KafkaConn = kafkaConn

	brokers, err := kafkaConn.Brokers()
	if err != nil {
		return errors.Wrap(err, "kafkaConn.Brokers")
	}

	s.Log.Infof("kafka connected to brokers: %+v", brokers)

	return nil
}

func (s *Server) initKafkaTopics(ctx context.Context) {
	controller, err := s.KafkaConn.Controller()
	if err != nil {
		s.Log.WarnMsg("kafkaConn.Controller", err)
		return
	}

	controllerURI := net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port))
	s.Log.Infof("kafka controller uri: %s", controllerURI)

	conn, err := kafka.DialContext(ctx, "tcp", controllerURI)
	if err != nil {
		s.Log.WarnMsg("initKafkaTopics.DialContext", err)
		return
	}
	defer conn.Close() // nolint: errcheck

	s.Log.Infof("established new kafka controller connection: %s", controllerURI)

	productCreateTopic := kafka.TopicConfig{
		Topic:             s.Cfg.KafkaTopics.ProductCreate.TopicName,
		NumPartitions:     s.Cfg.KafkaTopics.ProductCreate.Partitions,
		ReplicationFactor: s.Cfg.KafkaTopics.ProductCreate.ReplicationFactor,
	}

	productCreatedTopic := kafka.TopicConfig{
		Topic:             s.Cfg.KafkaTopics.ProductCreated.TopicName,
		NumPartitions:     s.Cfg.KafkaTopics.ProductCreated.Partitions,
		ReplicationFactor: s.Cfg.KafkaTopics.ProductCreated.ReplicationFactor,
	}

	productUpdateTopic := kafka.TopicConfig{
		Topic:             s.Cfg.KafkaTopics.ProductUpdate.TopicName,
		NumPartitions:     s.Cfg.KafkaTopics.ProductUpdate.Partitions,
		ReplicationFactor: s.Cfg.KafkaTopics.ProductUpdate.ReplicationFactor,
	}

	productUpdatedTopic := kafka.TopicConfig{
		Topic:             s.Cfg.KafkaTopics.ProductUpdated.TopicName,
		NumPartitions:     s.Cfg.KafkaTopics.ProductUpdated.Partitions,
		ReplicationFactor: s.Cfg.KafkaTopics.ProductUpdated.ReplicationFactor,
	}

	productDeleteTopic := kafka.TopicConfig{
		Topic:             s.Cfg.KafkaTopics.ProductDelete.TopicName,
		NumPartitions:     s.Cfg.KafkaTopics.ProductDelete.Partitions,
		ReplicationFactor: s.Cfg.KafkaTopics.ProductDelete.ReplicationFactor,
	}

	productDeletedTopic := kafka.TopicConfig{
		Topic:             s.Cfg.KafkaTopics.ProductDeleted.TopicName,
		NumPartitions:     s.Cfg.KafkaTopics.ProductDeleted.Partitions,
		ReplicationFactor: s.Cfg.KafkaTopics.ProductDeleted.ReplicationFactor,
	}

	if err := conn.CreateTopics(
		productCreateTopic,
		productUpdateTopic,
		productCreatedTopic,
		productUpdatedTopic,
		productDeleteTopic,
		productDeletedTopic,
	); err != nil {
		s.Log.WarnMsg("kafkaConn.CreateTopics", err)
		return
	}

	s.Log.Infof("kafka topics created or already exists: %+v", []kafka.TopicConfig{productCreateTopic, productUpdateTopic, productCreatedTopic, productUpdatedTopic, productDeleteTopic, productDeletedTopic})
}

func (s *Server) initMongoDBCollections(ctx context.Context) {
	err := s.MongoClient.Database(s.Cfg.Mongo.Db).CreateCollection(ctx, s.Cfg.MongoCollections.Products)
	if err != nil {
		if !utils.CheckErrMessages(err, catalog_constants.ErrMsgMongoCollectionAlreadyExists) {
			s.Log.Warnf("(CreateCollection) err: {%v}", err)
		}
	}

	indexOptions := options.Index().SetSparse(true).SetUnique(true)
	index, err := s.MongoClient.Database(s.Cfg.Mongo.Db).Collection(s.Cfg.MongoCollections.Products).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: consts.ProductIdIndex, Value: 1}},
		Options: indexOptions,
	})
	if err != nil && !utils.CheckErrMessages(err, catalog_constants.ErrMsgAlreadyExists) {
		s.Log.Warnf("(CreateOne) err: {%v}", err)
	}
	s.Log.Infof("(CreatedIndex) index: {%s}", index)

	list, err := s.MongoClient.Database(s.Cfg.Mongo.Db).Collection(s.Cfg.MongoCollections.Products).Indexes().List(ctx)
	if err != nil {
		s.Log.Warnf("(initMongoDBCollections) [List] err: {%v}", err)
	}

	if list != nil {
		var results []bson.M
		if err := list.All(ctx, &results); err != nil {
			s.Log.Warnf("(All) err: {%v}", err)
		}
		s.Log.Infof("(indexes) results: {%#v}", results)
	}

	collections, err := s.MongoClient.Database(s.Cfg.Mongo.Db).ListCollectionNames(ctx, bson.M{})
	if err != nil {
		s.Log.Warnf("(ListCollections) err: {%v}", err)
	}
	s.Log.Infof("(Collections) created collections: {%v}", collections)
}

func (s *Server) runMetrics(cancel context.CancelFunc) {
	metricsServer := echo.New()
	go func() {
		metricsServer.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
			StackSize:         catalog_constants.StackSize,
			DisablePrintStack: true,
			DisableStackAll:   true,
		}))
		metricsServer.GET(s.Cfg.Probes.PrometheusPath, echo.WrapHandler(promhttp.Handler()))
		s.Log.Infof("Metrics server is running on port: %s", s.Cfg.Probes.PrometheusPort)
		if err := metricsServer.Start(s.Cfg.Probes.PrometheusPort); err != nil {
			s.Log.Errorf("metricsServer.Start: %v", err)
			cancel()
		}
	}()
}

func (s *Server) getConsumerGroupTopics() []string {
	return []string{
		s.Cfg.KafkaTopics.ProductCreate.TopicName,
		s.Cfg.KafkaTopics.ProductUpdate.TopicName,
		s.Cfg.KafkaTopics.ProductDelete.TopicName,
	}
}
