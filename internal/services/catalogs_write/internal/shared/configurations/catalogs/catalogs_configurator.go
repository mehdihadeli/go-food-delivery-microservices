package catalogs

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/fxapp"
	customEcho "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/custom_echo"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/configurations"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/shared/configurations/catalogs/infrastructure"
)

type CatalogsConfigurator struct {
	*fxapp.Application
	infrastructureConfigurator *infrastructure.InfrastructureConfigurator
	productsModuleConfigurator *configurations.ProductsModuleConfigurator
}

func NewCatalogsConfigurator(fxapp *fxapp.Application) *CatalogsConfigurator {
	infraConfigurator := infrastructure.NewInfrastructureConfigurator(fxapp)
	productModuleConfigurator := configurations.NewProductsModuleConfigurator(fxapp)

	return &CatalogsConfigurator{
		Application:                fxapp,
		infrastructureConfigurator: infraConfigurator,
		productsModuleConfigurator: productModuleConfigurator,
	}
}

func (ic *CatalogsConfigurator) ConfigureCatalogs() {
	// Shared
	// Infrastructure
	ic.infrastructureConfigurator.ConfigInfrastructures()

	// Shared
	// Catalogs configurations
	ic.ResolveFunc(func(gorm *gorm.DB) error {
		err := ic.migrateCatalogs(gorm)
		if err != nil {
			return err
		}

		return nil
	})

	// Modules
	// Product module
	ic.productsModuleConfigurator.ConfigureProductsModule()
}

func (ic *CatalogsConfigurator) MapCatalogsEndpoints() {
	// Shared
	ic.ResolveFunc(
		func(catalogsServer customEcho.EchoHttpServer, options *config.AppOptions) error {
			catalogsServer.SetupDefaultMiddlewares()

			// Config catalogs root endpoint
			catalogsServer.RouteBuilder().
				RegisterRoutes(func(e *echo.Echo) {
					e.GET("", func(ec echo.Context) error {
						return ec.String(
							http.StatusOK,
							fmt.Sprintf("%s is running...", options.GetMicroserviceNameUpper()),
						)
					})
				})

			// Config catalogs swagger
			ic.configSwagger(catalogsServer.RouteBuilder())

			return nil
		},
	)

	// Modules
	// Products Module endpoints
	ic.productsModuleConfigurator.MapProductsEndpoints()
}
