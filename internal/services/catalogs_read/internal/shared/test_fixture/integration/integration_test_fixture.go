package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/bus"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/config"
	sharedContracts "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/shared/contracts"
)

const (
	DatabaseName          = "catalogs_write"
	ProductCollectionName = "products"
)

type IntegrationTestSharedFixture struct {
	Cfg *config.AppOptions
	Log logger.Logger
	suite.Suite
}

type IntegrationTestFixture struct {
	RedisProductRepository contracts.ProductCacheRepository
	MongoProductRepository contracts.ProductRepository
	Bus                    bus.Bus
	CatalogsMetrics        *sharedContracts.CatalogsMetrics
	Ctx                    context.Context
	cancel                 context.CancelFunc
}

func NewIntegrationTestSharedFixture(t *testing.T) *IntegrationTestSharedFixture {
	//	// we could use EmptyLogger if we don't want to log anything
	//	log := defaultLogger.Logger
	//	cfg, _ := config.NewConfig(constants.Test)
	//
	//	err := mappings.ConfigureProductsMappings()
	//	if err != nil {
	//		require.FailNow(t, err.Error())
	//	}
	//	require.NoError(t, err)
	//
	//	integration := &IntegrationTestSharedFixture{
	//		Cfg: cfg,
	//		Log: log,
	//	}
	//
	//	return integration
	//}
	//
	//func NewIntegrationTestFixture(shared *IntegrationTestSharedFixture) *IntegrationTestFixture {
	//	ctx, cancel := context.WithCancel(context.Background())
	//
	//	// we could use EmptyLogger if we don't want to log anything
	//	c := infrastructure.NewTestInfrastructureConfigurator(shared.T(), shared.Log, shared.Cfg)
	//	infrastructures, cleanup, err := c.ConfigInfrastructures(ctx)
	//	if err != nil {
	//		cancel()
	//		require.FailNow(shared.T(), err.Error())
	//	}
	//
	//	productRep := repositories.NewMongoProductRepository(
	//		infrastructures.Log,
	//		infrastructures.MongoClient,
	//	)
	//	redisRepository := repositories.NewRedisProductRepository(
	//		infrastructures.Log,
	//		infrastructures.Cfg,
	//		infrastructures.Redis,
	//	)
	//
	//	mqBus, err := rabbitmqTestContainer.NewRabbitMQTestContainers().
	//		Start(ctx, shared.T(), func(builder rabbitmqConfigurations.RabbitMQConfigurationBuilder) {
	//			// Products RabbitMQ configuration
	//			rabbitmq.ConfigProductsRabbitMQ(builder, infrastructures)
	//		})
	//	if err != nil {
	//		cancel()
	//		require.FailNow(shared.T(), err.Error())
	//	}
	//
	//	catalogsMetrics, err := metrics.ConfigCatalogsMetrics(
	//		infrastructures.Cfg,
	//		infrastructures.Metrics,
	//	)
	//	if err != nil {
	//		cancel()
	//		require.FailNow(shared.T(), err.Error())
	//	}
	//
	//	shared.T().Cleanup(func() {
	//		// with Cancel() we send signal to done() channel to stop  grpc, http and workers gracefully
	//		// https://dev.to/mcaci/how-to-use-the-context-done-method-in-go-22me
	//		// https://www.digitalocean.com/community/tutorials/how-to-use-contexts-in-go
	//		mediatr.ClearRequestRegistrations()
	//		cancel()
	//		cleanup()
	//	})
	//
	//	integration := &IntegrationTestFixture{
	//		InfrastructureConfigurations: infrastructures,
	//		Bus:                          mqBus,
	//		CatalogsMetrics:              catalogsMetrics,
	//		MongoProductRepository:       productRep,
	//		RedisProductRepository:       redisRepository,
	//		Ctx:                          ctx,
	//		cancel:                       cancel,
	//	}
	//
	//	return integration

	return &IntegrationTestSharedFixture{}
}

func (e *IntegrationTestFixture) Run() {
	// wait for consumers ready to consume before publishing messages, preparation background workers takes a bit time (for preventing messages lost)
	time.Sleep(1 * time.Second)
}
