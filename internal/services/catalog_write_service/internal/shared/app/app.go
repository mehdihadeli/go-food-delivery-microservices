package app

import (
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
	app.Run()
}
