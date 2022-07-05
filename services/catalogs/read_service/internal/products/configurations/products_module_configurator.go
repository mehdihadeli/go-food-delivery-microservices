package configurations

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/configurations"
)

type ProductsModuleConfigurator interface {
	ConfigureProductsModule() error
}

type productsModuleConfigurator struct {
	*configurations.Infrastructure
}

type ProductModule struct {
	*configurations.Infrastructure
	Mediator      *mediatr.Mediator
	ProductsGroup *echo.Group
}

func NewProductsModuleConfigurator(infrastructure *configurations.Infrastructure) *productsModuleConfigurator {
	return &productsModuleConfigurator{Infrastructure: infrastructure}
}

func (c *productsModuleConfigurator) ConfigureProductsModule(ctx context.Context) error {

	pm := &ProductModule{Infrastructure: c.Infrastructure}

	v1 := c.Echo.Group("/api/v1")
	pm.ProductsGroup = v1.Group("/" + c.Cfg.Http.ProductsPath)

	mediator, err := shared.NewCatalogsMediator(c.Log, c.Cfg, c.KafkaProducer)

	if err != nil {
		return err
	}

	pm.Mediator = mediator

	pm.configEndpoints(ctx)
	pm.configKafkaConsumers(ctx)

	if c.Cfg.DeliveryType == "grpc" {
		pm.configGrpc()
	}

	return nil
}
