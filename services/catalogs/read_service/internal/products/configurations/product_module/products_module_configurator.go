package product_module

import (
	"context"
	grpcServer "github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	customEcho "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/custom_echo"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/configurations/mappings"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/configurations/mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/data/repositories"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/configurations/infrastructure"
)

type productsModuleConfigurator struct {
	*infrastructure.InfrastructureConfigurations
	echoServer customEcho.EchoHttpServer
	grpcServer grpcServer.GrpcServer
}

func NewProductsModuleConfigurator(infrastructure *infrastructure.InfrastructureConfigurations, echoServer customEcho.EchoHttpServer, grpcServer grpcServer.GrpcServer) contracts.ProductsModuleConfigurator {
	return &productsModuleConfigurator{InfrastructureConfigurations: infrastructure, echoServer: echoServer, grpcServer: grpcServer}
}

func (c *productsModuleConfigurator) ConfigureProductsModule(ctx context.Context) error {

	mongoProductRepository := repositories.NewMongoProductRepository(c.Log, c.Cfg, c.MongoClient)
	redisRepository := repositories.NewRedisRepository(c.Log, c.Cfg, c.Redis)

	err := mappings.ConfigureMappings()
	if err != nil {
		return err
	}

	err = mediatr.ConfigProductsMediator(mongoProductRepository, redisRepository, c.InfrastructureConfigurations)
	if err != nil {
		return err
	}

	c.configKafkaConsumers(ctx)

	if c.Cfg.DeliveryType == "grpc" {
		c.configGrpc(ctx)
	} else {
		c.configEndpoints(ctx)
	}

	return nil
}
