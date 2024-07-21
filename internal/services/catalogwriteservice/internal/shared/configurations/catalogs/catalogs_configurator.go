package catalogs

import (
	"fmt"
	"net/http"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/config/environment"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/fxapp/contracts"
	echocontracts "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/customecho/contracts"
	migrationcontracts "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/migration/contracts"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/config"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/configurations"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/shared/configurations/catalogs/infrastructure"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type CatalogsServiceConfigurator struct {
	contracts.Application
	infrastructureConfigurator *infrastructure.InfrastructureConfigurator
	productsModuleConfigurator *configurations.ProductsModuleConfigurator
}

func NewCatalogsServiceConfigurator(
	app contracts.Application,
) *CatalogsServiceConfigurator {
	infraConfigurator := infrastructure.NewInfrastructureConfigurator(app)
	productModuleConfigurator := configurations.NewProductsModuleConfigurator(
		app,
	)

	return &CatalogsServiceConfigurator{
		Application:                app,
		infrastructureConfigurator: infraConfigurator,
		productsModuleConfigurator: productModuleConfigurator,
	}
}

func (ic *CatalogsServiceConfigurator) ConfigureCatalogs() error {
	// Shared
	// Infrastructure
	ic.infrastructureConfigurator.ConfigInfrastructures()

	// Shared
	// Catalogs configurations
	ic.ResolveFunc(
		func(db *gorm.DB, postgresMigrationRunner migrationcontracts.PostgresMigrationRunner) error {
			err := ic.migrateCatalogs(postgresMigrationRunner)
			if err != nil {
				return err
			}

			if ic.Environment() != environment.Test {
				err = ic.seedCatalogs(db)
				if err != nil {
					return err
				}
			}

			return nil
		},
	)

	// Modules
	// Product module
	err := ic.productsModuleConfigurator.ConfigureProductsModule()

	return err
}

func (ic *CatalogsServiceConfigurator) MapCatalogsEndpoints() error {
	// Shared
	ic.ResolveFunc(
		func(catalogsServer echocontracts.EchoHttpServer, options *config.AppOptions) error {
			catalogsServer.SetupDefaultMiddlewares()

			// config catalogs root endpoint
			catalogsServer.RouteBuilder().
				RegisterRoutes(func(e *echo.Echo) {
					e.GET("", func(ec echo.Context) error {
						return ec.String(
							http.StatusOK,
							fmt.Sprintf(
								"%s is running...",
								options.GetMicroserviceNameUpper(),
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
	err := ic.productsModuleConfigurator.MapProductsEndpoints()

	return err
}
