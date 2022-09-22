package integration

import (
	"context"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/contracts/store"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/eventstroredb"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger/defaultLogger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/bus"
	bus2 "github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/bus"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/configurations/mappings"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/contracts/repositories"
	orderRepositories "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/data/repositories"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/models/orders/aggregate"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/configurations/infrastructure"
	"math"
	"time"
)

type IntegrationTestFixture struct {
	*infrastructure.InfrastructureConfiguration
	OrderAggregateStore      store.AggregateStore[*aggregate.Order]
	MongoOrderReadRepository repositories.OrderReadRepository
	EsdbWorker               eventstroredb.EsdbSubscriptionAllWorker
	RabbitMQBus              bus.Bus
	ctx                      context.Context
	cancel                   context.CancelFunc
	Cleanup                  func()
}

func NewIntegrationTestFixture() *IntegrationTestFixture {
	deadline := time.Now().Add(time.Duration(math.MaxInt64))
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	cfg, _ := config.InitConfig("test")
	c := infrastructure.NewInfrastructureConfigurator(defaultLogger.Logger, cfg)
	infrastructures, _, cleanup := c.ConfigInfrastructures(ctx)

	eventStore := eventstroredb.NewEventStoreDbEventStore(infrastructures.Log, infrastructures.Esdb, infrastructures.EsdbSerializer)
	orderAggregateStore := eventstroredb.NewEventStoreAggregateStore[*aggregate.Order](infrastructures.Log, eventStore, infrastructures.EsdbSerializer)

	mongoOrderReadRepository := orderRepositories.NewMongoOrderReadRepository(infrastructures.Log, infrastructures.Cfg, infrastructures.MongoClient)

	esdbSubscribeAllWorker := eventstroredb.NewEsdbSubscriptionAllWorker(
		infrastructures.Log,
		infrastructures.Esdb,
		infrastructures.Cfg.EventStoreConfig,
		infrastructures.EsdbSerializer,
		infrastructures.CheckpointRepository,
		es.NewProjectionPublisher(infrastructures.Projections))

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
		InfrastructureConfiguration: infrastructures,
		OrderAggregateStore:         orderAggregateStore,
		MongoOrderReadRepository:    mongoOrderReadRepository,
		EsdbWorker:                  esdbSubscribeAllWorker,
		RabbitMQBus:                 rabbitmqBus,
		ctx:                         ctx,
		cancel:                      cancel,
	}
}

func (e *IntegrationTestFixture) Run() {
	go func() {
		//https://developers.eventstore.com/clients/grpc/subscriptions.html#filtering-by-prefix-1
		option := &eventstroredb.EventStoreDBSubscriptionToAllOptions{
			FilterOptions: &esdb.SubscriptionFilter{
				Type:     esdb.StreamFilterType,
				Prefixes: e.Cfg.Subscriptions.OrderSubscription.Prefix,
			},
			SubscriptionId: e.Cfg.Subscriptions.OrderSubscription.SubscriptionId,
		}

		err := e.EsdbWorker.SubscribeAll(e.ctx, option)
		if err != nil {
			e.cancel()
			e.Log.Errorf("(esdbSubscribeAllWorker.SubscribeAll) err: {%v}", err)
		}
	}()
	go func() {
		err := e.RabbitMQBus.Start(e.ctx)
		if err != nil {
			e.cancel()
			e.Log.Errorf("(RabbitMQBus.Start) err: {%v}", err)
		}
	}()
}
