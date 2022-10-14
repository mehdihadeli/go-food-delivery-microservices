package integration

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/constants"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger/defaultLogger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/bus"
	webWoker "github.com/mehdihadeli/store-golang-microservice-sample/pkg/web"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/configurations/mappings"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/data/repositories"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/configurations/catalogs/metrics"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/configurations/catalogs/rabbitmq"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/configurations/infrastructure"
	contracts2 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/web/workers"
)

type IntegrationTestFixture struct {
	*contracts2.InfrastructureConfigurations
	ProductRepository contracts.ProductRepository
	Bus               bus.Bus
	CatalogsMetrics   *contracts2.CatalogsMetrics
	workersRunner     *webWoker.WorkersRunner
	Ctx               context.Context
	cancel            context.CancelFunc
	Cleanup           func()
}

func NewIntegrationTestFixture() *IntegrationTestFixture {
	cfg, _ := config.InitConfig(constants.Test)

	ctx, cancel := context.WithCancel(context.Background())
	c := infrastructure.NewInfrastructureConfigurator(defaultLogger.Logger, cfg)
	infrastructures, cleanup, err := c.ConfigInfrastructures(context.Background())
	if err != nil {
		cancel()
		return nil
	}

	productRep := repositories.NewPostgresProductRepository(infrastructures.Log, cfg, infrastructures.Gorm.DB)

	err = mappings.ConfigureProductsMappings()
	if err != nil {
		cancel()
		return nil
	}

	mq, err := rabbitmq.ConfigCatalogsRabbitMQ(ctx, cfg.RabbitMQ, infrastructures)
	if err != nil {
		cancel()
		return nil
	}

	catalogsMetrics, err := metrics.ConfigCatalogsMetrics(cfg, infrastructures.Metrics)
	if err != nil {
		cancel()
		return nil
	}

	workersRunner := webWoker.NewWorkersRunner([]webWoker.Worker{
		workers.NewRabbitMQWorker(infrastructures.Log, mq),
	})

	return &IntegrationTestFixture{
		Cleanup: func() {
			workersRunner.Stop(ctx)
			cancel()
			cleanup()
		},
		InfrastructureConfigurations: infrastructures,
		Bus:                          mq,
		CatalogsMetrics:              catalogsMetrics,
		ProductRepository:            productRep,
		workersRunner:                workersRunner,
		Ctx:                          ctx,
		cancel:                       cancel,
	}
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
}
