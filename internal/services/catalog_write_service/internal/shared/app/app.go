package app

import (
	"context"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/shared/configurations/catalogs"
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
	app.ConfigureCatalogs()

	app.MapCatalogsEndpoints()

	app.Logger().Info("Starting catalog_service application")
	app.ResolveFunc(func(tracer tracing.AppTracer) {
		_, span := tracer.Start(context.Background(), "Application started")
		span.End()
	})

	app.Run()
}
