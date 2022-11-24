package integration

import (
	"context"
	"math"
	"time"

	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/constants"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/contracts/store"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/eventstroredb"
	defaultLogger2 "github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger/default_logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/bus"
	webWoker "github.com/mehdihadeli/store-golang-microservice-sample/pkg/web"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/configurations/mappings"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/contracts/repositories"
	orderRepositories "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/data/repositories"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/models/orders/aggregate"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/configurations/infrastructure"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/configurations/orders/metrics"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/configurations/orders/rabbitmq"
	subscriptionAll "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/configurations/orders/subscription_all"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/web/workers"
)

type IntegrationTestFixture struct {
	contracts.InfrastructureConfigurations
	OrderAggregateStore      store.AggregateStore[*aggregate.Order]
	MongoOrderReadRepository repositories.OrderReadRepository
	Bus                      bus.Bus
	OrdersMetrics            contracts.OrdersMetrics
	workersRunner            *webWoker.WorkersRunner
	Ctx                      context.Context
	cancel                   context.CancelFunc
	Cleanup                  func()
	cleanupChan              chan struct{}
}

func NewIntegrationTestFixture() *IntegrationTestFixture {
	cfg, _ := config.InitConfig(constants.Test)

	deadline := time.Now().Add(time.Duration(math.MaxInt64))
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	c := infrastructure.NewInfrastructureConfigurator(defaultLogger2.Logger, cfg)
	infrastructures, cleanup, err := c.ConfigInfrastructures(ctx)
	if err != nil {
		cancel()
		return nil
	}

	eventStore := eventstroredb.NewEventStoreDbEventStore(infrastructures.Log(), infrastructures.Esdb(), infrastructures.EsdbSerializer())
	orderAggregateStore := eventstroredb.NewEventStoreAggregateStore[*aggregate.Order](infrastructures.Log(), eventStore, infrastructures.EsdbSerializer())

	mongoOrderReadRepository := orderRepositories.NewMongoOrderReadRepository(infrastructures.Log(), infrastructures.Cfg(), infrastructures.MongoClient())

	err = mappings.ConfigureOrdersMappings()
	if err != nil {
		cancel()
		return nil
	}

	mq, err := rabbitmq.ConfigOrdersRabbitMQ(ctx, cfg.RabbitMQ, infrastructures)
	if err != nil {
		cancel()
		return nil
	}

	if err != nil {
		cancel()
		return nil
	}

	subscriptionAllWorker, err := subscriptionAll.ConfigOrdersSubscriptionAllWorker(infrastructures, mq)
	if err != nil {
		cancel()
		return nil
	}

	ordersMetrics, err := metrics.ConfigOrdersMetrics(cfg, infrastructures.Metrics())
	if err != nil {
		cancel()
		return nil
	}

	workersRunner := webWoker.NewWorkersRunner([]webWoker.Worker{
		workers.NewRabbitMQWorker(infrastructures.Log(), mq), workers.NewEventStoreDBWorker(infrastructures.Log(), infrastructures.Cfg(), subscriptionAllWorker),
	})

	return &IntegrationTestFixture{
		Cleanup: func() {
			cancel()
			cleanup()
		},
		workersRunner:                workersRunner,
		Bus:                          mq,
		OrdersMetrics:                ordersMetrics,
		InfrastructureConfigurations: infrastructures,
		OrderAggregateStore:          orderAggregateStore,
		MongoOrderReadRepository:     mongoOrderReadRepository,
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

	go func() {
		for {
			select {
			case _ = <-workersErr:
				e.cancel()
				return
				// case <-e.cleanupChan:
				//	workersRunner.Stop(e.Ctx)
				//	return
			}
		}
	}()
}
