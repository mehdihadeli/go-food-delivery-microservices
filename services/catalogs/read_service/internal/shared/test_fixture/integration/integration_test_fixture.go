package integration

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/constants"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger/defaultLogger"
	webWoker "github.com/mehdihadeli/store-golang-microservice-sample/pkg/web"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/configurations/consumers"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/configurations/mappings"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/data/repositories"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/configurations/infrastructure"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/web/workers"
)

type IntegrationTestFixture struct {
	*infrastructure.InfrastructureConfigurations
	RedisProductRepository contracts.ProductCacheRepository
	MongoProductRepository contracts.ProductRepository
	workersRunner          *webWoker.WorkersRunner
	ctx                    context.Context
	cancel                 context.CancelFunc
	Cleanup                func()
}

func NewIntegrationTestFixture() *IntegrationTestFixture {
	cfg, _ := config.InitConfig(constants.Test)

	ctx, cancel := context.WithCancel(context.Background())
	c := infrastructure.NewInfrastructureConfigurator(defaultLogger.Logger, cfg)
	infrastructures, _, cleanup := c.ConfigInfrastructures(context.Background())

	mongoProductRepository := repositories.NewMongoProductRepository(infrastructures.Log, cfg, infrastructures.MongoClient)
	redisProductRepository := repositories.NewRedisRepository(infrastructures.Log, cfg, infrastructures.Redis)

	err := mappings.ConfigureMappings()
	if err != nil {
		cancel()
		return nil
	}

	err = consumers.ConfigConsumers(infrastructures)
	if err != nil {
		cancel()
		return nil
	}

	workersRunner := webWoker.NewWorkersRunner([]webWoker.Worker{
		workers.NewRabbitMQWorkerWorker(infrastructures),
	})

	return &IntegrationTestFixture{
		Cleanup: func() {
			workersRunner.Stop(ctx)
			cancel()
			cleanup()
		},
		InfrastructureConfigurations: infrastructures,
		RedisProductRepository:       redisProductRepository,
		MongoProductRepository:       mongoProductRepository,
		workersRunner:                workersRunner,
		ctx:                          ctx,
		cancel:                       cancel,
	}
}

func (e *IntegrationTestFixture) Run() {
	workersErr := e.workersRunner.Start(e.ctx)
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
