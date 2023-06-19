package main

import (
	"context"

	"go.uber.org/fx"

	application "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/shared/app"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/shared/configurations/catalogs"
)

// https://github.com/swaggo/swag#how-to-use-it-with-gin

// @contact.name Mehdi Hadeli
// @contact.url https://github.com/mehdihadeli
// @title Catalogs Write-Service Api
// @version 1.0
// @description Catalogs Write-Service Api.
func main() {
	// configure dependencies
	appBuilder := application.NewCatalogsWriteApplicationBuilder()
	appBuilder.ProvideModule(catalogs.CatalogsServiceModule)

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
	app.ConfigureCatalogs()

	app.MapCatalogsEndpoints()

	app.Logger.Info("Starting catalog_service application")
	app.Run()
}
