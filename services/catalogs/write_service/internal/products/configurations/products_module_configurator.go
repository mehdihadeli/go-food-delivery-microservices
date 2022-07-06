package configurations

import (
	"context"
	repositories_imp "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/data/repositories"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/configurations/infrastructure"
)

type ProductsModuleConfigurator interface {
	ConfigureProductsModule() error
}

type productsModuleConfigurator struct {
	*infrastructure.InfrastructureConfiguration
}

func NewProductsModuleConfigurator(infrastructure *infrastructure.InfrastructureConfiguration) *productsModuleConfigurator {
	return &productsModuleConfigurator{InfrastructureConfiguration: infrastructure}
}

func (c *productsModuleConfigurator) ConfigureProductsModule(ctx context.Context) error {

	v1 := c.Echo.Group("/api/v1")
	group := v1.Group("/" + c.Cfg.Http.ProductsPath)

	productRepository := repositories_imp.NewPostgresProductRepository(c.Log, c.Cfg, c.Gorm)

	mediator, err := c.configProductsMediator(productRepository)

	if err != nil {
		return err
	}

	c.configEndpoints(ctx, group, mediator)

	if c.Cfg.DeliveryType == "grpc" {
		c.configGrpc(ctx, mediator)
	}

	return nil
}
