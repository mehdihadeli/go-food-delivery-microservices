package e2e

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/constants"
	grpcServer "github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	customEcho "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/custom_echo"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	defaultLogger "github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger/default_logger"
	mongoRepository "github.com/mehdihadeli/store-golang-microservice-sample/pkg/mongodb/repository"
	rabbitmqConfigurations "github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/configurations"
	rabbitmqTestContainer "github.com/mehdihadeli/store-golang-microservice-sample/pkg/test/containers/testcontainer/rabbitmq"
	webWorker "github.com/mehdihadeli/store-golang-microservice-sample/pkg/web"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/configurations/mappings"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/configurations/mediator"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/configurations/rabbitmq"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/data/repositories"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/delivery"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/models"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/configurations/catalogs/metrics"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/configurations/infrastructure"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/web/workers"
)

const (
	DatabaseName          = "catalogs"
	ProductCollectionName = "products"
)

type E2ETestSharedFixture struct {
	Cfg *config.Config
	Log logger.Logger
	suite.Suite
}

type E2ETestFixture struct {
	GrpcServer grpcServer.GrpcServer
	HttpServer customEcho.EchoHttpServer
	*delivery.ProductEndpointBase
	workersRunner *webWorker.WorkersRunner
	Ctx           context.Context
	cancel        context.CancelFunc
}

func NewE2ETestSharedFixture(t *testing.T) *E2ETestSharedFixture {
	// we could use EmptyLogger if we don't want to log anything
	log := defaultLogger.Logger
	cfg, _ := config.InitConfig(constants.Test)

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
	ctx, cancel := context.WithCancel(context.Background())

	c := infrastructure.NewTestInfrastructureConfigurator(shared.T(), shared.Log, shared.Cfg)
	infrastructures, cleanup, err := c.ConfigInfrastructures(ctx)
	if err != nil {
		cancel()
		return nil
	}

	genericRepo := mongoRepository.NewGenericMongoRepository[*models.Product](infrastructures.MongoClient, DatabaseName, ProductCollectionName)
	productRep := repositories.NewMongoProductRepository(infrastructures.Log, genericRepo)
	redisRepository := repositories.NewRedisProductRepository(infrastructures.Log, infrastructures.Cfg, infrastructures.Redis)

	mqBus, err := rabbitmqTestContainer.NewRabbitMQTestContainers().Start(ctx, shared.T(), func(builder rabbitmqConfigurations.RabbitMQConfigurationBuilder) {
		// Products RabbitMQ configuration
		rabbitmq.ConfigProductsRabbitMQ(builder, infrastructures)
	})
	if err != nil {
		cancel()
		require.FailNow(shared.T(), err.Error())
	}

	catalogsMetrics, err := metrics.ConfigCatalogsMetrics(infrastructures.Cfg, infrastructures.Metrics)
	if err != nil {
		cancel()
		require.FailNow(shared.T(), err.Error())
	}

	err = mediator.ConfigProductsMediator(infrastructures, productRep, redisRepository, mqBus)
	if err != nil {
		cancel()
		require.FailNow(shared.T(), err.Error())
	}

	grpcServer := grpcServer.NewGrpcServer(infrastructures.Cfg.GRPC, defaultLogger.Logger, infrastructures.Cfg.ServiceName, infrastructures.Metrics)
	httpServer := customEcho.NewEchoHttpServer(infrastructures.Cfg.Http, defaultLogger.Logger, infrastructures.Cfg.ServiceName, infrastructures.Metrics)
	httpServer.SetupDefaultMiddlewares()

	var productEndpointBase *delivery.ProductEndpointBase

	httpServer.RouteBuilder().RegisterGroupFunc("/api/v1", func(v1 *echo.Group) {
		group := v1.Group("/products")
		productEndpointBase = delivery.NewProductEndpointBase(infrastructures, group, mqBus, catalogsMetrics)
	})

	httpServer.RouteBuilder().RegisterRoutes(func(e *echo.Echo) {
		e.GET("", func(ec echo.Context) error {
			return ec.String(http.StatusOK, fmt.Sprintf("%s is running...", infrastructures.Cfg.GetMicroserviceNameUpper()))
		})
	})

	workersRunner := webWorker.NewWorkersRunner([]webWorker.Worker{
		workers.NewRabbitMQWorker(infrastructures.Log, mqBus),
	})

	shared.T().Cleanup(func() {
		// with Cancel() we send signal to done() channel to stop grpc, http and workers gracefully
		// https://dev.to/mcaci/how-to-use-the-context-done-method-in-go-22me
		// https://www.digitalocean.com/community/tutorials/how-to-use-contexts-in-go
		mediatr.ClearRequestRegistrations()
		cancel()
		cleanup()
	})

	return &E2ETestFixture{
		GrpcServer:          grpcServer,
		HttpServer:          httpServer,
		ProductEndpointBase: productEndpointBase,
		workersRunner:       workersRunner,
		Ctx:                 ctx,
		cancel:              cancel,
	}
}

func (e *E2ETestFixture) Run() {
	go func() {
		if err := e.GrpcServer.RunGrpcServer(e.Ctx, nil); err != nil {
			e.Log.Errorf("(s.RunGrpcServer) err: %v", err)
		}
	}()

	go func() {
		if err := e.HttpServer.RunHttpServer(e.Ctx, nil); err != nil {
			e.Log.Errorf("(s.RunHttpServer) err: %v", err)
		}
	}()

	workersErr := e.workersRunner.Start(e.Ctx)
	go func() {
		for {
			select {
			case _ = <-workersErr:
				e.cancel()
				return
			}
		}
	}()

	// wait for consumers ready to consume before publishing messages, preparation background workers takes a bit time (for preventing messages lost)
	time.Sleep(1 * time.Second)
}
