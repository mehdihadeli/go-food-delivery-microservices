package product_module

import (
	"context"
	grpcServer "github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	customEcho "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/custom_echo"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/configurations/mappings"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/configurations/mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/contracts"
	contracts2 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/contracts"
)

type productsModuleConfigurator struct {
	contracts2.InfrastructureConfiguration
	echoServer customEcho.EchoHttpServer
	grpcServer grpcServer.GrpcServer
}

func NewProductsModuleConfigurator(infrastructure contracts2.InfrastructureConfiguration, echoServer customEcho.EchoHttpServer, grpcServer grpcServer.GrpcServer) contracts.ProductsModuleConfigurator {
	return &productsModuleConfigurator{InfrastructureConfiguration: infrastructure, echoServer: echoServer, grpcServer: grpcServer}
}

func (c *productsModuleConfigurator) ConfigureProductsModule(ctx context.Context) error {
	err := mappings.ConfigureMappings()
	if err != nil {
		return err
	}

	err = mediatr.ConfigProductsMediator(c.InfrastructureConfiguration)
	if err != nil {
		return err
	}

	if c.GetCfg().DeliveryType == "grpc" {
		c.configGrpc(ctx)
	} else {
		c.configEndpoints(ctx)
	}

	return nil
}
