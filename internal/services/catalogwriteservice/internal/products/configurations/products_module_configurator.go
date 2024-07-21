package configurations

import (
	fxcontracts "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/fxapp/contracts"
	grpcServer "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/grpc"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/configurations/endpoints"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/configurations/mappings"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/configurations/mediator"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/shared/grpc"
	productsservice "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/shared/grpc/genproto"

	googleGrpc "google.golang.org/grpc"
)

type ProductsModuleConfigurator struct {
	fxcontracts.Application
}

func NewProductsModuleConfigurator(
	fxapp fxcontracts.Application,
) *ProductsModuleConfigurator {
	return &ProductsModuleConfigurator{
		Application: fxapp,
	}
}

func (c *ProductsModuleConfigurator) ConfigureProductsModule() error {
	// config products mappings
	err := mappings.ConfigureProductsMappings()
	if err != nil {
		return err
	}

	// register products request handler on mediator
	c.ResolveFuncWithParamTag(
		mediator.RegisterMediatorHandlers,
		`group:"product-handlers"`,
	)

	return nil
}

func (c *ProductsModuleConfigurator) MapProductsEndpoints() error {
	// config endpoints
	c.ResolveFuncWithParamTag(
		endpoints.RegisterEndpoints,
		`group:"product-routes"`,
	)

	// config Products Grpc Endpoints
	c.ResolveFunc(
		func(catalogsGrpcServer grpcServer.GrpcServer, grpcService *grpc.ProductGrpcServiceServer) error {
			catalogsGrpcServer.GrpcServiceBuilder().
				RegisterRoutes(func(server *googleGrpc.Server) {
					productsservice.RegisterProductsServiceServer(
						server,
						grpcService,
					)
				})

			return nil
		},
	)

	return nil
}
