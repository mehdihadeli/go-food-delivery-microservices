package app

import "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/orderservice/internal/shared/configurations/orders"

type App struct{}

func NewApp() *App {
	return &App{}
}

func (a *App) Run() {
	// configure dependencies
	appBuilder := NewOrdersApplicationBuilder()
	appBuilder.ProvideModule(orders.OrderServiceModule)

	app := appBuilder.Build()

	// configure application
	app.ConfigureOrders()

	app.MapOrdersEndpoints()

	app.Logger().Info("Starting orders_service application")
	app.Run()
}
