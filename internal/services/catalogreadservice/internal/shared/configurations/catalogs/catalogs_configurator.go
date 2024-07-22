package catalogs

import (
	"fmt"
	"net/http"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/fxapp/contracts"
	echocontracts "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/customecho/contracts"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogreadservice/config"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogreadservice/internal/products/configurations"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogreadservice/internal/shared/configurations/catalogs/infrastructure"

	"github.com/labstack/echo/v4"
)

type CatalogsServiceConfigurator struct {
	contracts.Application
	infrastructureConfigurator *infrastructure.InfrastructureConfigurator
	productsModuleConfigurator *configurations.ProductsModuleConfigurator
}

func NewCatalogsServiceConfigurator(app contracts.Application) *CatalogsServiceConfigurator {
	infraConfigurator := infrastructure.NewInfrastructureConfigurator(app)
	productModuleConfigurator := configurations.NewProductsModuleConfigurator(app)

	return &CatalogsServiceConfigurator{
		Application:                app,
		infrastructureConfigurator: infraConfigurator,
		productsModuleConfigurator: productModuleConfigurator,
	}
}

func (ic *CatalogsServiceConfigurator) ConfigureCatalogs() {
	// Shared
	// Infrastructure
	ic.infrastructureConfigurator.ConfigInfrastructures()

	// Shared
	// Catalogs configurations

	// Modules
	// Product module
	ic.productsModuleConfigurator.ConfigureProductsModule()
}

func (ic *CatalogsServiceConfigurator) MapCatalogsEndpoints() {
	// Shared
	ic.ResolveFunc(
		func(catalogsServer echocontracts.EchoHttpServer, cfg *config.Config) error {
			catalogsServer.SetupDefaultMiddlewares()

			// config catalogs root endpoint
			catalogsServer.RouteBuilder().
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

			// config catalogs swagger
			ic.configSwagger(catalogsServer.RouteBuilder())

			return nil
		},
	)

	// Modules
	// Products CatalogsServiceModule endpoints
	ic.productsModuleConfigurator.MapProductsEndpoints()
}
