package e2e

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/constants"
	grpcServer "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/grpc"
	customEcho "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/custom_echo"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	defaultLogger "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/default_logger"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/products/configurations/mappings"
)

type E2ETestSharedFixture struct {
	Cfg *config.AppOptions
	Log logger.Logger
	suite.Suite
}

type E2ETestFixture struct {
	GrpcServer grpcServer.GrpcServer
	HttpServer customEcho.EchoHttpServer
	Ctx        context.Context
	cancel     context.CancelFunc
}

func NewE2ETestSharedFixture(t *testing.T) *E2ETestSharedFixture {
	// we could use EmptyLogger if we don't want to log anything
	log := defaultLogger.Logger
	cfg, _ := config.NewConfig(constants.Test)

	err := mappings.ConfigureProductsMappings()
	if err != nil {
		require.FailNow(t, err.Error())
	}
	require.NoError(t, err)

	integration := &E2ETestSharedFixture{
		Cfg: cfg,
		Log: log,
	}

	return integration
}

func NewE2ETestFixture(shared *E2ETestSharedFixture) *E2ETestFixture {
	//ctx, cancel := context.WithCancel(context.Background())
	//
	//c := infrastructure.NewTestInfrastructureConfigurator(shared.T(), shared.Log, shared.Cfg)
	//infrastructures, cleanup, err := c.ConfigInfrastructures(ctx)
	//if err != nil {
	//	cancel()
	//	return nil
	//}
	//
	//productRep := repositories.NewMongoProductRepository(
	//	infrastructures.Log,
	//	infrastructures.MongoClient,
	//)
	//redisRepository := repositories.NewRedisProductRepository(
	//	infrastructures.Log,
	//	infrastructures.Cfg,
	//	infrastructures.Redis,
	//)
	//
	//mqBus, err := rabbitmqTestContainer.NewRabbitMQTestContainers().
	//	Start(ctx, shared.T(), func(builder rabbitmqConfigurations.RabbitMQConfigurationBuilder) {
	//		// Products RabbitMQ configuration
	//		rabbitmq.ConfigProductsRabbitMQ(builder, infrastructures)
	//	})
	//if err != nil {
	//	cancel()
	//	require.FailNow(shared.T(), err.Error())
	//}
	//
	//catalogsMetrics, err := metrics.ConfigCatalogsMetrics(
	//	infrastructures.Cfg,
	//	infrastructures.Metrics,
	//)
	//if err != nil {
	//	cancel()
	//	require.FailNow(shared.T(), err.Error())
	//}
	//
	//err = mediator.ConfigProductsMediator(infrastructures, productRep, redisRepository, mqBus)
	//if err != nil {
	//	cancel()
	//	require.FailNow(shared.T(), err.Error())
	//}
	//
	//grpcOptions, _ := config2.BindConfigKey[*config3.GrpcOptions](
	//	strcase.ToLowerCamel(typeMapper.GetTypeNameByT[config3.GrpcOptions]()),
	//	constants.Test,
	//)
	//
	//echoOptions, _ := config2.BindConfigKey[*config4.EchoHttpOptions](
	//	strcase.ToLowerCamel(typeMapper.GetTypeNameByT[config4.EchoHttpOptions]()),
	//	constants.Test,
	//)
	//
	//grpcServer := grpcServer.NewGrpcServer(
	//	grpcOptions,
	//	defaultLogger.Logger,
	//	infrastructures.Metrics,
	//)
	//httpServer := customEcho.NewEchoHttpServer(customEcho.EchoHttpServerParams{
	//	Logger: defaultLogger.Logger,
	//	Config: echoOptions,
	//	Meter:  infrastructures.Metrics,
	//},
	//)
	//httpServer.SetupDefaultMiddlewares()
	//
	//var productEndpointBase *delivery.ProductEndpointBase
	//
	//httpServer.RouteBuilder().RegisterGroupFunc("/api/v1", func(v1 *echo.Group) {
	//	group := v1.Group("/products")
	//	productEndpointBase = delivery.NewProductEndpointBase(
	//		infrastructures,
	//		group,
	//		mqBus,
	//		catalogsMetrics,
	//	)
	//})
	//
	//httpServer.RouteBuilder().RegisterRoutes(func(e *echo.Echo) {
	//	e.GET("", func(ec echo.Context) error {
	//		return ec.String(
	//			http.StatusOK,
	//			fmt.Sprintf("%s is running...", infrastructures.Cfg.GetMicroserviceNameUpper()),
	//		)
	//	})
	//})
	//
	//shared.T().Cleanup(func() {
	//	// with Cancel() we send signal to done() channel to stop grpc, http and workers gracefully
	//	// https://dev.to/mcaci/how-to-use-the-context-done-method-in-go-22me
	//	// https://www.digitalocean.com/community/tutorials/how-to-use-contexts-in-go
	//	mediatr.ClearRequestRegistrations()
	//	cancel()
	//	cleanup()
	//})
	//
	//return &E2ETestFixture{
	//	GrpcServer:          grpcServer,
	//	HttpServer:          httpServer,
	//	ProductEndpointBase: productEndpointBase,
	//	Ctx:                 ctx,
	//	cancel:              cancel,
	//}

	return &E2ETestFixture{}
}

func (e *E2ETestFixture) Run() {
	// wait for consumers ready to consume before publishing messages, preparation background workers takes a bit time (for preventing messages lost)
	time.Sleep(1 * time.Second)
}
