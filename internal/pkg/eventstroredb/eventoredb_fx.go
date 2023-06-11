package eventstroredb

import (
	"context"
	"time"

	"github.com/EventStore/EventStore-Client-Go/esdb"
	"go.uber.org/fx"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/eventstroredb/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
)

// Module provided to fxlog
// https://uber-go.github.io/fx/modules.html
var Module = func(builder ProjectionBuilderFuc) fx.Option {
	return fx.Module(
		"eventstoredbfx",
		// - order is not important in provide
		// - provide can have parameter and will resolve if registered
		// - execute its func only if it requested
		fx.Provide(
			config.ProvideConfig,
			NewEsdbSerializer,
			NewEventStoreDB,
			NewEventStoreDbEventStore,
			NewEsdbSubscriptionCheckpointRepository,
			NewEsdbSubscriptionAllWorker,
		),
		fx.Supply(builder),
		// - execute after registering all of our provided
		// - they execute by their orders
		// - invokes always execute its func compare to provides that only run when we request for them.
		// - return value will be discarded and can not be provided
		fx.Invoke(registerHooks),
	)
}

// we don't want to register any dependencies here, its func body should execute always even we don't request for that, so we should use `invoke`
func registerHooks(
	lc fx.Lifecycle,
	worker EsdbSubscriptionAllWorker,
	logger logger.Logger,
	cfg config.EventStoreDbOptions,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Start server in a separate goroutine, this way when the server is shutdown "s.e.Start" will
			// return promptly, and the call to "s.e.Shutdown" is the one that will wait for all other
			// resources to be properly freed. If it was the other way around, the application would just
			// exit without gracefully shutting down the server.
			// For more details: https://medium.com/@momchil.dev/proper-http-shutdown-in-go-bd3bfaade0f2
			go func() {
				option := &EventStoreDBSubscriptionToAllOptions{
					FilterOptions: &esdb.SubscriptionFilter{
						Type:     esdb.StreamFilterType,
						Prefixes: cfg.Subscription.Prefix,
					},
					SubscriptionId: cfg.Subscription.SubscriptionId,
				}

				if err := worker.SubscribeAll(ctx, option); err != nil {
					logger.Fatalf(
						"(EsdbSubscriptionAllWorker.Start) error in running esdb subscription worker: {%v}",
						err,
					)
				}
			}()
			logger.Info("esdb subscription worker is listening.")

			return nil
		},
		OnStop: func(ctx context.Context) error {
			ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
			defer cancel()

			return nil
		},
	})
}
