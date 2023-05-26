package workers

import (
	"context"

	"github.com/EventStore/EventStore-Client-Go/esdb"

	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/eventstroredb"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/web"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/config"
)

func NewEventStoreDBWorker(
	logger logger.Logger,
	cfg *config.Config,
	subscriptionAllWorker eventstroredb.EsdbSubscriptionAllWorker,
) web.Worker {
	return web.NewBackgroundWorker(func(ctx context.Context) error {
		option := &eventstroredb.EventStoreDBSubscriptionToAllOptions{
			FilterOptions: &esdb.SubscriptionFilter{
				Type:     esdb.StreamFilterType,
				Prefixes: cfg.Subscriptions.OrderSubscription.Prefix,
			},
			SubscriptionId: cfg.Subscriptions.OrderSubscription.SubscriptionId,
		}
		err := subscriptionAllWorker.SubscribeAll(ctx, option)
		if err != nil {
			logger.Errorf(
				"[EventStoreDBWorker.SubscribeAll] error in the subscribing eventstore: {%v}",
				err,
			)
			return err
		}
		return nil
	}, nil)
}
