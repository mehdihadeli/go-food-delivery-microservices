package configurations

import (
	"github.com/go-playground/validator"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/fxapp"
	grpcServer "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/grpc"
	customEcho "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/custom_echo"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/producer"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/configurations/endpoints"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/configurations/grpc"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/configurations/mappings"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/configurations/mediatr"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/contracts/data"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/shared/contracts"
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
		func(logger logger.Logger, uow data.CatalogUnitOfWork, productRepository data.ProductRepository, producer producer.Producer) error {
			// Config Products Mediators
			err := mediatr.ConfigProductsMediator(logger, uow, productRepository, producer)
			if err != nil {
				return err
			}

			// cfg Products Mappings
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
		func(logger logger.Logger, validator *validator.Validate, catalogsMetrics *contracts.CatalogsMetrics, catalogsServer customEcho.EchoHttpServer, catalogsGrpcServer grpcServer.GrpcServer) error {
			// Config Http endpoints
			endpoints.ConfigProductsEndpoints(
				catalogsServer.RouteBuilder(),
				catalogsMetrics,
				validator,
				logger,
			)

			// Config Products gRPC endpoints
			grpc.ConfigProductsGrpc(
				catalogsGrpcServer.GrpcServiceBuilder(),
				logger,
				catalogsMetrics,
			)

			return nil
		},
	)
}
