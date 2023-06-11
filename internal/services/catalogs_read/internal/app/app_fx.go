package app

import (
	"context"
	"fmt"

	"go.uber.org/fx"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/bus"
	appconfig "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/shared/configurations/infrastructure"
)

var Module = fx.Module("appfx",
	// infrastructure setup --> should go infrastructure module
	infrastructure.Module,

	// application setup
	appconfig.Module,

	// https://uber-go.github.io/fx/lifecycle.html#lifecycle-hooks
	fx.Invoke(registerHooks),
	fx.Invoke(func(bus bus.RabbitmqBus) {
		fmt.Print("Creating")
	}),
)

func registerHooks(lifecycle fx.Lifecycle) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}
