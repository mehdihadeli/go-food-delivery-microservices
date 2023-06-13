package rabbitmq

import (
	"context"
	"time"

	"go.uber.org/fx"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	bus2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/bus"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/producer"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/bus"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/config"
)

// Module provided to fxlog
// https://uber-go.github.io/fx/modules.html
var Module = func(rabbitMQConfigurationConstructor interface{}) fx.Option {
	return fx.Module(
		"rabbitmqfx",
		// - order is not important in provide
		// - provide can have parameter and will resolve if registered
		// - execute its func only if it requested
		fx.Provide(
			config.ProvideConfig,
		),
		fx.Provide(fx.Annotate(
			bus.NewRabbitmqBus,
			fx.ParamTags(``, ``, ``, `optional:"true"`),
		)),
		fx.Provide(fx.Annotate(
			bus.NewRabbitmqBus,
			fx.ParamTags(``, ``, ``, `optional:"true"`),
			fx.As(new(bus2.Bus)),
		)),
		fx.Provide(fx.Annotate(
			bus.NewRabbitmqBus,
			fx.ParamTags(``, ``, ``, `optional:"true"`),
			fx.As(new(producer.Producer)),
		)),
		fx.Provide(rabbitMQConfigurationConstructor),
		//// without return type
		//// fxlog.Invoke(rabbitmqBuilderFunc),
		//// https://github.com/uber-go/fx/pull/833
		//// https://pkg.go.dev/go.uber.org/fx#Decorate
		//fxapp.Decorate(
		//	json.NewJsonSerializer,
		//	serializer.NewDefaultEventSerializer,
		//	serializer.NewDefaultMetadataSerializer,
		//),
		//// https://github.com/uber-go/fx/pull/837
		//// https://pkg.go.dev/go.uber.org/fx#Replace
		//fxapp.Replace(zap.Module),

		// - execute after registering all of our provided
		// - they execute by their orders
		// - invokes always execute its func compare to provides that only run when we request for them.
		// - return value will be discarded and can not be provided
		fx.Invoke(registerHooks),
	)
}

// we don't want to register any dependencies here, its func body should execute always even we don't request for that, so we should use `invoke`
func registerHooks(lc fx.Lifecycle, bus bus.RabbitmqBus, logger logger.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Start server in a separate goroutine, this way when the server is shutdown "s.e.Start" will
			// return promptly, and the call to "s.e.Shutdown" is the one that will wait for all other
			// resources to be properly freed. If it was the other way around, the application would just
			// exit without gracefully shutting down the server.
			// For more details: https://medium.com/@momchil.dev/proper-http-shutdown-in-go-bd3bfaade0f2
			go func() {
				if err := bus.Start(ctx); err != nil {
					logger.Fatalf("(bus.Start) error in running rabbitmq server: {%v}", err)
				}
			}()
			logger.Info("rabbitmq is listening.")

			return nil
		},
		OnStop: func(ctx context.Context) error {
			ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
			defer cancel()
			if err := bus.Stop(ctx); err != nil {
				logger.Errorf("error shutting down rabbitmq server: %v", err)
			} else {
				logger.Info("rabbitmq server shutdown gracefully")
			}
			return nil
		},
	})
}
