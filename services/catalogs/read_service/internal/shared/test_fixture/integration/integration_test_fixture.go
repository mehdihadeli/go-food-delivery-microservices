package integration

import (
	"context"
	"testing"
	"time"

	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/configurations/rabbitmq"

	"github.com/mehdihadeli/go-mediatr"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/constants"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	defaultLogger "github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger/default_logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/bus"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mongodb/repository"
	rabbitmqConfigurations "github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/configurations"
	rabbitmqTestContainer "github.com/mehdihadeli/store-golang-microservice-sample/pkg/test/containers/testcontainer/rabbitmq"
	webWorker "github.com/mehdihadeli/store-golang-microservice-sample/pkg/web"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/configurations/mappings"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/data/repositories"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/models"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/configurations/catalogs/metrics"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/configurations/infrastructure"
	sharedContracts "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/web/workers"
)

const (
	DatabaseName          = "catalogs"
	ProductCollectionName = "products"
)

type IntegrationTestSharedFixture struct {
	Cfg *config.Config
	Log logger.Logger
	suite.Suite
}

type IntegrationTestFixture struct {
	*sharedContracts.InfrastructureConfigurations
	RedisProductRepository contracts.ProductCacheRepository
	MongoProductRepository contracts.ProductRepository
	Bus                    bus.Bus
	CatalogsMetrics        *sharedContracts.CatalogsMetrics
	workersRunner          *webWorker.WorkersRunner
	Ctx                    context.Context
	cancel                 context.CancelFunc
}

func NewIntegrationTestSharedFixture(t *testing.T) *IntegrationTestSharedFixture {
	// we could use EmptyLogger if we don't want to log anything
	log := defaultLogger.Logger
	cfg, _ := config.InitConfig(constants.Test)

	err := mappings.ConfigureProductsMappings()
	if err != nil {
		require.FailNow(t, err.Error())
	}
	require.NoError(t, err)

	integration := &IntegrationTestSharedFixture{
		Cfg: cfg,
		Log: log,
	}

	return integration
}

func NewIntegrationTestFixture(shared *IntegrationTestSharedFixture) *IntegrationTestFixture {
	ctx, cancel := context.WithCancel(context.Background())

	// we could use EmptyLogger if we don't want to log anything
	c := infrastructure.NewTestInfrastructureConfigurator(shared.T(), shared.Log, shared.Cfg)
	infrastructures, cleanup, err := c.ConfigInfrastructures(ctx)
	if err != nil {
		cancel()
		require.FailNow(shared.T(), err.Error())
	}

	genericRepo := repository.NewGenericMongoRepository[*models.Product](infrastructures.MongoClient, DatabaseName, ProductCollectionName)
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

	workersRunner := webWorker.NewWorkersRunner([]webWorker.Worker{
		workers.NewRabbitMQWorker(infrastructures.Log, mqBus),
	})

	shared.T().Cleanup(func() {
		// with Cancel() we send signal to done() channel to stop  grpc, http and workers gracefully
		// https://dev.to/mcaci/how-to-use-the-context-done-method-in-go-22me
		// https://www.digitalocean.com/community/tutorials/how-to-use-contexts-in-go
		mediatr.ClearRequestRegistrations()
		cancel()
		cleanup()
	})

	integration := &IntegrationTestFixture{
		InfrastructureConfigurations: infrastructures,
		Bus:                          mqBus,
		CatalogsMetrics:              catalogsMetrics,
		MongoProductRepository:       productRep,
		RedisProductRepository:       redisRepository,
		workersRunner:                workersRunner,
		Ctx:                          ctx,
		cancel:                       cancel,
	}

	return integration
}

func (e *IntegrationTestFixture) Run() {
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
