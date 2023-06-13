package configurations

import (
	"github.com/go-playground/validator"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/fxapp"
	customEcho "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/custom_echo"
	logger2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	bus2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/bus"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/products/configurations/endpoints"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/products/configurations/mappings"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/products/configurations/mediator"
	contracts2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/products/contracts"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/shared/contracts"
)

type ProductsModuleConfigurator struct {
	*fxapp.Application
}

func NewProductsModuleConfigurator(
	fxapp *fxapp.Application,
) *ProductsModuleConfigurator {
	return &ProductsModuleConfigurator{
		Application: fxapp,
	}
}

func (c *ProductsModuleConfigurator) ConfigureProductsModule() {
	c.ResolveFunc(
		func(logger logger2.Logger, mongoRepository contracts2.ProductRepository, cacheRepository contracts2.ProductCacheRepository, bus bus2.Bus) error {
			// Config Products Mediators
			err := mediator.ConfigProductsMediator(logger, mongoRepository, cacheRepository, bus)
			if err != nil {
				return err
			}

			// Config Products Mappings
			err = mappings.ConfigureProductsMappings()
			if err != nil {
				return err
			}
			return nil
		},
	)
}

func (c *ProductsModuleConfigurator) MapProductsEndpoints() {
	c.ResolveFunc(
		// Config Products Endpoints
		func(logger logger2.Logger, validator *validator.Validate, catalogsMetrics *contracts.CatalogsMetrics, catalogsServer customEcho.EchoHttpServer) error {
			endpoints.ConfigProductsEndpoints(
				catalogsServer.RouteBuilder(),
				catalogsMetrics,
				validator,
				logger,
			)

			return nil
		},
	)
}
