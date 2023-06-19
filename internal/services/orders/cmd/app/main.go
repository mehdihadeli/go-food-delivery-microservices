package main

import (
	"context"

	application "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orders/internal/shared/app"
	orders "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orders/internal/shared/configurations/orders"

	"go.uber.org/fx"
)

// https://github.com/swaggo/swag#how-to-use-it-with-gin

// @contact.name Mehdi Hadeli
// @contact.url https://github.com/mehdihadeli
// @title Orders Service Api
// @version 1.0
// @description Orders Service Api
func main() {
	// configure dependencies
	appBuilder := application.NewOrdersApplicationBuilder()
	appBuilder.ProvideModule(orders.OrderServiceModule)

	app := appBuilder.Build()

	app.RegisterHook(func(lifecycle fx.Lifecycle) {
		lifecycle.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				return nil
			},
			OnStop: func(ctx context.Context) error {
				// some cleanup if exists
				return nil
			},
		})
	})

	// configure application
	app.ConfigureOrders()

	app.MapOrdersEndpoints()

	app.Logger.Info("Starting orders_service application")
	app.Run()
}
