package integration

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/constants"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger/defaultLogger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/bus"
	bus2 "github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/bus"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/configurations/mappings"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/data/repositories"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/configurations/infrastructure"
)

type IntegrationTestFixture struct {
	*infrastructure.InfrastructureConfigurations
	RedisProductRepository contracts.ProductCacheRepository
	MongoProductRepository contracts.ProductRepository
	RabbitMQBus            bus.Bus
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

	rabbitmqBus := bus2.NewRabbitMQBus(infrastructures.Log, infrastructures.Consumers)

	err := mappings.ConfigureMappings()
	if err != nil {
		cancel()
		return nil
	}

	return &IntegrationTestFixture{
		Cleanup: func() {
			cancel()
			cleanup()
			rabbitmqBus.Stop(ctx)
		},
		InfrastructureConfigurations: infrastructures,
		RedisProductRepository:       redisProductRepository,
		MongoProductRepository:       mongoProductRepository,
		RabbitMQBus:                  rabbitmqBus,
		ctx:                          ctx,
		cancel:                       cancel,
	}
}

func (e *IntegrationTestFixture) Run() {
	go func() {
		err := e.RabbitMQBus.Start(e.ctx)
		if err != nil {
			e.cancel()
			e.Log.Errorf("(RabbitMQBus.Start) err: {%v}", err)
		}
	}()
}
