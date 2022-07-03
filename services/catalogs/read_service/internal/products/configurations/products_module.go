package configurations

import (
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/configurations"
)

type ProductsModuleConfigurator interface {
	ConfigureProductsModule() error
}

type productsModuleConfigurator struct {
	infrastructure *configurations.Infrastructure
}

type ProductModule struct {
	Infrastructure *configurations.Infrastructure
	Mediator       *mediatr.Mediator
	ProductsGroup  *echo.Group
}

func NewProductsModuleConfigurator(infrastructure *configurations.Infrastructure) *productsModuleConfigurator {
	return &productsModuleConfigurator{infrastructure: infrastructure}
}

func (c *productsModuleConfigurator) ConfigureProductsModule() error {

	pm := ProductModule{Infrastructure: c.infrastructure}

	v1 := c.infrastructure.Echo.Group("/api/v1")
	pm.ProductsGroup = v1.Group("/" + c.infrastructure.Cfg.Http.ProductsPath)

	m, err := shared.NewCatalogsMediator(c.infrastructure.Log, c.infrastructure.Cfg, c.infrastructure.KafkaProducer)

	if err != nil {
		return err
	}

	pm.Mediator = m

	pm.configEndpoints()

	if c.infrastructure.Cfg.DeliveryType == "grpc" {
		pm.configGrpc()
	}

	return nil
}

func (pm *ProductModule) configEndpoints() {

}

func (pm *ProductModule) configGrpc() {
	//productGrpcService := grpc.NewProductGrpcService(pm.Infrastructure, pm.Mediator, pm.ProductRepository)
	//product_service.RegisterProductsServiceServer(pm.Infrastructure.GrpcServer, productGrpcService)
}
