package configurations

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/data/repositories"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/mappings"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/configurations/infrastructure"
)

type ProductsModuleConfigurator interface {
	ConfigureProductsModule() error
}

type productsModuleConfigurator struct {
	*infrastructure.InfrastructureConfigurations
}

func NewProductsModuleConfigurator(infrastructure *infrastructure.InfrastructureConfigurations) *productsModuleConfigurator {
	return &productsModuleConfigurator{InfrastructureConfigurations: infrastructure}
}

func (c *productsModuleConfigurator) ConfigureProductsModule(ctx context.Context) error {

	v1 := c.Echo.Group("/api/v1")
	group := v1.Group("/" + c.Cfg.Http.ProductsPath)

	productRepository := repositories.NewMongoProductRepository(c.Log, c.Cfg, c.MongoClient)

	err := mappings.ConfigureMappings()
	if err != nil {
		return err
	}

	mediator, err := c.configProductsMediator(productRepository)

	if err != nil {
		return err
	}

	c.configEndpoints(ctx, mediator, group)
	c.configKafkaConsumers(ctx, mediator)

	if c.Cfg.DeliveryType == "grpc" {
		c.configGrpc(ctx, mediator)
	}

	return nil
}
