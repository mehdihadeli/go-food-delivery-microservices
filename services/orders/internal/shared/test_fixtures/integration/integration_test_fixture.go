package integration

import (
	"context"
	"emperror.dev/errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/constants"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/contracts/store"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/eventstroredb"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger/defaultLogger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/consumer/options"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test/messaging/consumer"
	webWoker "github.com/mehdihadeli/store-golang-microservice-sample/pkg/web"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/configurations/mappings"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/configurations/projections"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/contracts/repositories"
	orderRepositories "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/data/repositories"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/models/orders/aggregate"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/configurations/infrastructure"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/web/workers"
	"math"
	"time"
)

type IntegrationTestFixture struct {
	*infrastructure.InfrastructureConfiguration
	OrderAggregateStore      store.AggregateStore[*aggregate.Order]
	MongoOrderReadRepository repositories.OrderReadRepository
	ctx                      context.Context
	cancel                   context.CancelFunc
	Cleanup                  func()
	cleanupChan              chan struct{}
}

func NewIntegrationTestFixture() *IntegrationTestFixture {
	cfg, _ := config.InitConfig(constants.Test)

	deadline := time.Now().Add(time.Duration(math.MaxInt64))
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	c := infrastructure.NewInfrastructureConfigurator(defaultLogger.Logger, cfg)
	infrastructures, _, cleanup := c.ConfigInfrastructures(ctx)

	eventStore := eventstroredb.NewEventStoreDbEventStore(infrastructures.Log, infrastructures.Esdb, infrastructures.EsdbSerializer)
	orderAggregateStore := eventstroredb.NewEventStoreAggregateStore[*aggregate.Order](infrastructures.Log, eventStore, infrastructures.EsdbSerializer)

	mongoOrderReadRepository := orderRepositories.NewMongoOrderReadRepository(infrastructures.Log, infrastructures.Cfg, infrastructures.MongoClient)

	err := mappings.ConfigureMappings()
	if err != nil {
		cancel()
		return nil
	}

	projections.ConfigOrderProjections(infrastructures)

	cleanupChan := make(chan struct{})
	return &IntegrationTestFixture{
		cleanupChan: cleanupChan,
		Cleanup: func() {
			cleanupChan <- struct{}{}
			cancel()
			cleanup()
		},
		InfrastructureConfiguration: infrastructures,
		OrderAggregateStore:         orderAggregateStore,
		MongoOrderReadRepository:    mongoOrderReadRepository,
		ctx:                         ctx,
		cancel:                      cancel,
	}
}

func (e *IntegrationTestFixture) Run() {
	workersRunner := webWoker.NewWorkersRunner([]webWoker.Worker{
		workers.NewRabbitMQWorkerWorker(e.InfrastructureConfiguration), workers.NewEventStoreDBWorker(e.InfrastructureConfiguration),
	})

	workersErr := workersRunner.Start(e.ctx)
	go func() {
		for {
			select {
			case _ = <-workersErr:
				workersRunner.Stop(e.ctx)
				e.cancel()
				return
			case <-e.cleanupChan:
				workersRunner.Stop(e.ctx)
				return
			}
		}
	}()
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
