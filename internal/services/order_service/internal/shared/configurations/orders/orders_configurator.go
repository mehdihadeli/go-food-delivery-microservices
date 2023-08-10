package orders

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/fxapp/contracts"
	customEcho "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/custom_echo"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/configurations"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/shared/configurations/orders/infrastructure"
)

type OrdersServiceConfigurator struct {
	contracts.Application
	infrastructureConfigurator *infrastructure.InfrastructureConfigurator
	ordersModuleConfigurator   *configurations.OrdersModuleConfigurator
}

func NewOrdersServiceConfigurator(
	app contracts.Application,
) *OrdersServiceConfigurator {
	infraConfigurator := infrastructure.NewInfrastructureConfigurator(app)
	ordersModuleConfigurator := configurations.NewOrdersModuleConfigurator(app)

	return &OrdersServiceConfigurator{
		Application:                app,
		infrastructureConfigurator: infraConfigurator,
		ordersModuleConfigurator:   ordersModuleConfigurator,
	}
}

func (ic *OrdersServiceConfigurator) ConfigureOrders() {
	// Shared
	// Infrastructure
	ic.infrastructureConfigurator.ConfigInfrastructures()

	// Shared
	// Orders service configurations

	// Modules
	// Order module
	ic.ordersModuleConfigurator.ConfigureOrdersModule()
}

func (ic *OrdersServiceConfigurator) MapOrdersEndpoints() {
	// Shared
	ic.ResolveFunc(
		func(ordersServer customEcho.EchoHttpServer, cfg *config.Config) error {
			ordersServer.SetupDefaultMiddlewares()

			// Config orders root endpoint
			ordersServer.RouteBuilder().
				RegisterRoutes(func(e *echo.Echo) {
					e.GET("", func(ec echo.Context) error {
						return ec.String(
							http.StatusOK,
							fmt.Sprintf(
								"%s is running...",
								cfg.AppOptions.GetMicroserviceNameUpper(),
							),
						)
					})
				})

			// Config orders swagger
			ic.configSwagger(ordersServer.RouteBuilder())

			return nil
		},
	)

	// Modules
	// Orders Module endpoints
	ic.ordersModuleConfigurator.MapOrdersEndpoints()
}
