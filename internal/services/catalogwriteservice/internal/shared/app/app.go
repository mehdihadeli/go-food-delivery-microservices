package app

import (
	"context"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/otel/tracing"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/shared/configurations/catalogs"
)

type App struct{}

func NewApp() *App {
	return &App{}
}

func (a *App) Run() {
	// configure dependencies
	appBuilder := NewCatalogsWriteApplicationBuilder()
	appBuilder.ProvideModule(catalogs.CatalogsServiceModule)

	app := appBuilder.Build()

	// configure application
	err := app.ConfigureCatalogs()
	if err != nil {
		app.Logger().Fatalf("Error in ConfigureCatalogs", err)
	}

	err = app.MapCatalogsEndpoints()
	if err != nil {
		app.Logger().Fatalf("Error in MapCatalogsEndpoints", err)
	}

	app.Logger().Info("Starting catalog_service application")
	app.ResolveFunc(func(tracer tracing.AppTracer) {
		_, span := tracer.Start(context.Background(), "Application started")
		span.End()
	})

	app.Run()
}
