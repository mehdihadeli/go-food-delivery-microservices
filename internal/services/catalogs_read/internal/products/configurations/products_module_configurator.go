package configurations

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/fxapp"
	logger2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/web/route"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/products/configurations/mappings"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/products/configurations/mediator"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/products/contracts/data"
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
		func(logger logger2.Logger, mongoRepository data.ProductRepository, cacheRepository data.ProductCacheRepository, tracer tracing.AppTracer) error {
			// Config Products Mediators
			err := mediator.ConfigProductsMediator(logger, mongoRepository, cacheRepository, tracer)
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
	// Config Products Http Endpoints
	c.ResolveFuncWithParamTag(func(endpoints []route.Endpoint) {
		for _, endpoint := range endpoints {
			endpoint.MapEndpoint()
		}
	}, `group:"product-routes"`,
	)
}