package catalogs

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/fxapp"
	customEcho "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/custom_echo"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/products/configurations"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/shared/configurations/catalogs/infrastructure"
)

type CatalogsServiceConfigurator struct {
	*fxapp.Application
	infrastructureConfigurator *infrastructure.InfrastructureConfigurator
	productsModuleConfigurator *configurations.ProductsModuleConfigurator
}

func NewCatalogsServiceConfigurator(fxapp *fxapp.Application) *CatalogsServiceConfigurator {
	infraConfigurator := infrastructure.NewInfrastructureConfigurator(fxapp)
	productModuleConfigurator := configurations.NewProductsModuleConfigurator(fxapp)

	return &CatalogsServiceConfigurator{
		Application:                fxapp,
		infrastructureConfigurator: infraConfigurator,
		productsModuleConfigurator: productModuleConfigurator,
	}
}

func (ic *CatalogsServiceConfigurator) ConfigureCatalogs() {
	// Shared
	ic.infrastructureConfigurator.ConfigInfrastructures()

	// Modules
	// Product module
	ic.productsModuleConfigurator.ConfigureProductsModule()

	//catalogsServiceConfigurations := &contracts.CatalogsServiceConfigurations{}
	//
	//catalogsServiceConfigurations.CatalogsEchoServer.SetupDefaultMiddlewares()
	//
	//catalogsServiceConfigurations.CatalogsEchoServer.RouteBuilder().
	//	RegisterRoutes(func(e *echo.Echo) {
	//		e.GET("", func(ec echo.Context) error {
	//			return ec.String(
	//				http.StatusOK,
	//				fmt.Sprintf("%s is running...", ic.Cfg.GetMicroserviceNameUpper()),
	//			)
	//		})
	//	})
	//
	//// Catalogs Swagger Configs
	//ic.configSwagger(catalogsServiceConfigurations.CatalogsEchoServer.RouteBuilder())
	//
	//// Catalogs Metrics Configs
	//catalogsMetrics, err := catalogsMetrics.ConfigCatalogsMetrics(ic.Cfg, ic.Metrics)
	//if err != nil {
	//	return nil, err
	//}
	//catalogsServiceConfigurations.CatalogsMetrics = catalogsMetrics
	//
	//// Catalogs Product Module Configs
	//pc := configurations.NewProductsModuleConfigurator(
	//	ic.InfrastructureConfigurations,
	//	catalogsMetrics,
	//	bus,
	//	catalogsServiceConfigurations.CatalogsEchoServer.RouteBuilder(),
	//	catalogsServiceConfigurations.CatalogsGrpcServer.GrpcServiceBuilder(),
	//)
	//err = pc.ConfigureProductsModule(ctx)
	//if err != nil {
	//	return nil, errors.WithMessage(
	//		err,
	//		"[CatalogsServiceConfigurator_ConfigureCatalogsService.ConfigureProductsModule] error in product module configurator",
	//	)
	//}
}

func (ic *CatalogsServiceConfigurator) MapCatalogsEndpoints() {
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
