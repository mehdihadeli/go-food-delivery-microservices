package integration

import (
	"context"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/es/contracts/store"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/bus"
	webWoker "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/web"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orders/internal/orders/contracts/repositories"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orders/internal/orders/models/orders/aggregate"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orders/internal/shared/contracts"
)

type IntegrationTestFixture struct {
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
	//cfg, _ := config.NewConfig(constants.Test)
	//
	//deadline := time.Now().Add(time.Duration(math.MaxInt64))
	//ctx, cancel := context.WithDeadline(context.Background(), deadline)
	//c := infrastructure.NewInfrastructureConfigurator(defaultLogger2.Logger, cfg)
	//infrastructures, cleanup, err := c.ConfigInfrastructures(ctx)
	//if err != nil {
	//	cancel()
	//	return nil
	//}
	//
	//eventStore := eventstroredb.NewEventStoreDbEventStore(
	//	infrastructures.Log(),
	//	infrastructures.Esdb(),
	//	infrastructures.EsdbSerializer(),
	//)
	//orderAggregateStore := eventstroredb.NewEventStoreAggregateStore[*aggregate.Order](
	//	infrastructures.Log(),
	//	eventStore,
	//	infrastructures.EsdbSerializer(),
	//)
	//
	//mongoOrderReadRepository := orderRepositories.NewMongoOrderReadRepository(
	//	infrastructures.Log(),
	//	infrastructures.Cfg(),
	//	infrastructures.MongoClient(),
	//)
	//
	//err = mappings.ConfigureOrdersMappings()
	//if err != nil {
	//	cancel()
	//	return nil
	//}
	//
	//mq, err := rabbitmq.ConfigOrdersRabbitMQ(ctx, cfg.RabbitMQ, infrastructures)
	//if err != nil {
	//	cancel()
	//	return nil
	//}
	//
	//if err != nil {
	//	cancel()
	//	return nil
	//}
	//
	//subscriptionAllWorker, err := subscriptionAll.ConfigOrdersSubscriptionAllWorker(
	//	infrastructures,
	//	mq,
	//)
	//if err != nil {
	//	cancel()
	//	return nil
	//}
	//
	//ordersMetrics, err := metrics.ConfigOrdersMetrics(cfg, infrastructures.Metrics())
	//if err != nil {
	//	cancel()
	//	return nil
	//}
	//
	//workersRunner := webWoker.NewWorkersRunner([]webWoker.Worker{
	//	workers.NewRabbitMQWorker(
	//		infrastructures.Log(),
	//		mq,
	//	), workers.NewEventStoreDBWorker(infrastructures.Log(), infrastructures.Cfg(), subscriptionAllWorker),
	//})
	//
	//return &IntegrationTestFixture{
	//	Cleanup: func() {
	//		cancel()
	//		cleanup()
	//	},
	//	workersRunner:                workersRunner,
	//	Bus:                          mq,
	//	OrdersMetrics:                ordersMetrics,
	//	InfrastructureConfigurations: infrastructures,
	//	OrderAggregateStore:          orderAggregateStore,
	//	MongoOrderReadRepository:     mongoOrderReadRepository,
	//	Ctx:                          ctx,
	//	cancel:                       cancel,
	//}
	return &IntegrationTestFixture{}
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
