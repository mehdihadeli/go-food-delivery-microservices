package configurations

import (
	"context"
	grpcServer "github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	customEcho "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/custom_echo"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/configurations/endpoints"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/configurations/grpc"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/configurations/mappings"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/configurations/mediator"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/contracts"
	contracts2 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/contracts"
)

type productsModuleConfigurator struct {
	contracts2.InfrastructureConfigurations
	routeBuilder       *customEcho.RouteBuilder
	grpcServiceBuilder *grpcServer.GrpcServiceBuilder
}

func NewProductsModuleConfigurator(infrastructure contracts2.InfrastructureConfigurations, routeBuilder *customEcho.RouteBuilder, grpcServiceBuilder *grpcServer.GrpcServiceBuilder) contracts.ProductsModuleConfigurator {
	return &productsModuleConfigurator{InfrastructureConfigurations: infrastructure, routeBuilder: routeBuilder, grpcServiceBuilder: grpcServiceBuilder}
}

func (c *productsModuleConfigurator) ConfigureProductsModule(ctx context.Context) error {
	if c.Cfg().DeliveryType == "grpc" {
		//Config Products Grpc
		grpc.ConfigProductsGrpc(ctx, c.grpcServiceBuilder, c.InfrastructureConfigurations)
	} else {
		//Config Products Endpoints
		endpoints.ConfigProductsEndpoints(ctx, c.routeBuilder, c.InfrastructureConfigurations)
	}

	//Config Products Mappings
	err := mappings.ConfigeProductsMappings()
	if err != nil {
		return err
	}

	//Config Products Mediators
	err = mediator.ConfigProductsMediator(c.InfrastructureConfigurations)
	if err != nil {
		return err
	}

	return nil
}
