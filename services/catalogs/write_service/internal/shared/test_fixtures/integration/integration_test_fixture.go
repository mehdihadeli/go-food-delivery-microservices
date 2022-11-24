package integration

import (
	"context"
	"testing"
	"time"

	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/gorm_postgres/repository"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/models"

	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/data/uow"

	"github.com/mehdihadeli/go-mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/constants"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	defaultLogger "github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger/default_logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/bus"
	rabbitmqConfigurations "github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/configurations"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test/containers/testcontainer/rabbitmq"
	webWoker "github.com/mehdihadeli/store-golang-microservice-sample/pkg/web"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/configurations/mappings"
	rabbitmq2 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/configurations/rabbitmq"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts/data"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/data/repositories"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/configurations/catalogs/metrics"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/configurations/infrastructure"
	contracts2 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/web/workers"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type IntegrationTestSharedFixture struct {
	Cfg *config.Config
	Log logger.Logger
	suite.Suite
}

type IntegrationTestFixture struct {
	*contracts2.InfrastructureConfigurations
	ProductRepository  data.ProductRepository
	CatalogUnitOfWorks data.CatalogUnitOfWork
	Bus                bus.Bus
	CatalogsMetrics    *contracts2.CatalogsMetrics
	workersRunner      *webWoker.WorkersRunner
	Ctx                context.Context
	cancel             context.CancelFunc
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

	genericRepo := repository.NewGenericGormRepository[*models.Product](infrastructures.Gorm)
	productRep := repositories.NewPostgresProductRepository(infrastructures.Log, genericRepo)
	catalogUnitOfWork := uow.NewCatalogsUnitOfWork(infrastructures.Log, infrastructures.Gorm)

	mqBus, err := rabbitmq.NewRabbitMQTestContainers().Start(ctx, shared.T(), func(builder rabbitmqConfigurations.RabbitMQConfigurationBuilder) {
		// Products RabbitMQ configuration
		rabbitmq2.ConfigProductsRabbitMQ(builder)
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

	workersRunner := webWoker.NewWorkersRunner([]webWoker.Worker{
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
		ProductRepository:            productRep,
		CatalogUnitOfWorks:           catalogUnitOfWork,
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
