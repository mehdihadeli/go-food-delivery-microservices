package integration

import (
	"context"
	"emperror.dev/errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/constants"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger/defaultLogger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/consumer/options"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test/messaging/consumer"
	webWoker "github.com/mehdihadeli/store-golang-microservice-sample/pkg/web"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/configurations/mappings"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/data/repositories"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/configurations/infrastructure"
	"time"
)

type IntegrationTestFixture struct {
	*infrastructure.InfrastructureConfiguration
	ProductRepository contracts.ProductRepository
	workersRunner     *webWoker.WorkersRunner
	Ctx               context.Context
	cancel            context.CancelFunc
	Cleanup           func()
}

func NewIntegrationTestFixture() *IntegrationTestFixture {
	cfg, _ := config.InitConfig(constants.Test)

	ctx, cancel := context.WithCancel(context.Background())
	c := infrastructure.NewInfrastructureConfigurator(defaultLogger.Logger, cfg)
	infrastructures, _, cleanup := c.ConfigInfrastructures(context.Background())

	productRep := repositories.NewPostgresProductRepository(infrastructures.Log, cfg, infrastructures.Gorm.DB)

	err := mappings.ConfigureMappings()
	if err != nil {
		cancel()
		return nil
	}

	return &IntegrationTestFixture{
		Cleanup: func() {
			cancel()
			cleanup()
		},
		InfrastructureConfiguration: infrastructures,
		ProductRepository:           productRep,
		Ctx:                         ctx,
		cancel:                      cancel,
	}
}

func (e *IntegrationTestFixture) Run() {
}

func (e *IntegrationTestFixture) FakeConsumer(messageName string) *consumer.RabbitMQFakeTestConsumer {
	fakeConsumer := consumer.NewRabbitMQFakeTestConsumer(
		e.EventSerializer,
		e.Log,
		e.RabbitMQConnection,
		func(builder *options.RabbitMQConsumerOptionsBuilder) {
			builder.WithExchangeName(messageName).WithQueueName(messageName).WithRoutingKey(messageName)
		})

	e.Consumers = append(e.Consumers, fakeConsumer)

	return fakeConsumer
}

func (e *IntegrationTestFixture) WaitUntilConditionMet(conditionToMet func() bool) error {
	timeout := 20 * time.Second

	startTime := time.Now()
	timeOutExpired := false
	meet := conditionToMet()
	for meet == false {
		if timeOutExpired {
			return errors.New("Condition not met for the test, timeout exceeded")
		}
		time.Sleep(time.Second * 2)
		meet = conditionToMet()
		timeOutExpired = time.Now().Sub(startTime) > timeout
	}

	return nil
}
