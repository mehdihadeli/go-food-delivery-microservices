package orders

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/fxapp"
	customEcho "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/custom_echo"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orders/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orders/internal/orders/configurations"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orders/internal/shared/configurations/orders/infrastructure"
)

type OrdersServiceConfigurator struct {
	*fxapp.Application
	infrastructureConfigurator *infrastructure.InfrastructureConfigurator
	ordersModuleConfigurator   *configurations.OrdersModuleConfigurator
}

func NewOrdersServiceConfigurator(
	fxapp *fxapp.Application,
) *OrdersServiceConfigurator {
	infraConfigurator := infrastructure.NewInfrastructureConfigurator(fxapp)
	ordersModuleConfigurator := configurations.NewOrdersModuleConfigurator(fxapp)

	return &OrdersServiceConfigurator{
		Application:                fxapp,
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
