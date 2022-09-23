package workers

import (
	"context"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/eventstroredb"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/web"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/configurations/infrastructure"
)

func NewEventStoreDBWorker(infra *infrastructure.InfrastructureConfiguration) web.Worker {
	esdbWorker := eventstroredb.NewEsdbSubscriptionAllWorker(
		infra.Log,
		infra.Esdb,
		infra.Cfg.EventStoreConfig,
		infra.EsdbSerializer,
		infra.CheckpointRepository,
		es.NewProjectionPublisher(infra.Projections))

	return web.NewBackgroundWorker(func(ctx context.Context) error {
		option := &eventstroredb.EventStoreDBSubscriptionToAllOptions{
			FilterOptions: &esdb.SubscriptionFilter{
				Type:     esdb.StreamFilterType,
				Prefixes: infra.Cfg.Subscriptions.OrderSubscription.Prefix,
			},
			SubscriptionId: infra.Cfg.Subscriptions.OrderSubscription.SubscriptionId,
		}
		err := esdbWorker.SubscribeAll(ctx, option)
		if err != nil {
			infra.Log.Errorf("[EventStoreDBWorker.SubscribeAll] error in the subscribing eventstore: {%v}", err)
			return err
		}
		return nil
	}, nil)

}
