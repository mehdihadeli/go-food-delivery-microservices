package product_mdoule

import (
	"context"
	grpcServer "github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	customEcho "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/custom_echo"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/configurations/mappings"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/configurations/mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts"
	repositoriesImp "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/data/repositories"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/configurations/infrastructure"
)

type productsModuleConfigurator struct {
	*infrastructure.InfrastructureConfiguration
	echoServer customEcho.EchoHttpServer
	grpcServer grpcServer.GrpcServer
}

func NewProductsModuleConfigurator(infrastructure *infrastructure.InfrastructureConfiguration, echoServer customEcho.EchoHttpServer, grpcServer grpcServer.GrpcServer) contracts.ProductsModuleConfigurator {
	return &productsModuleConfigurator{InfrastructureConfiguration: infrastructure, echoServer: echoServer, grpcServer: grpcServer}
}

func (c *productsModuleConfigurator) ConfigureProductsModule(ctx context.Context) error {
	productRepository := repositoriesImp.NewPostgresProductRepository(c.Log, c.Cfg, c.Gorm.DB)

	err := mappings.ConfigureMappings()
	if err != nil {
		return err
	}

	err = mediatr.ConfigProductsMediator(productRepository, c.InfrastructureConfiguration)
	if err != nil {
		return err
	}

	if c.Cfg.DeliveryType == "grpc" {
		c.configGrpc(ctx)
	} else {
		c.configEndpoints(ctx)
	}

	return nil
}
