package configurations

import (
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts"
	repositories_imp "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/data/repositories"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/configurations"
)

type ProductModule struct {
	*configurations.Infrastructure
	Mediator          *mediatr.Mediator
	ProductRepository contracts.ProductRepository
	ProductsGroup     *echo.Group
}

type ProductsModuleConfigurator interface {
	ConfigureProductsModule() error
}

type productsModuleConfigurator struct {
	*configurations.Infrastructure
}

func NewProductsModuleConfigurator(infrastructure *configurations.Infrastructure) *productsModuleConfigurator {
	return &productsModuleConfigurator{Infrastructure: infrastructure}
}

func (c *productsModuleConfigurator) ConfigureProductsModule() error {

	pm := ProductModule{Infrastructure: c.Infrastructure}

	v1 := c.Echo.Group("/api/v1")
	pm.ProductsGroup = v1.Group("/" + c.Cfg.Http.ProductsPath)

	pm.ProductRepository = repositories_imp.NewPostgresProductRepository(c.Log, c.Cfg, c.PgConn, c.Gorm)
	m, err := shared.NewCatalogsMediator(c.Log, c.Cfg, pm.ProductRepository, c.KafkaProducer)

	if err != nil {
		return err
	}

	pm.Mediator = m

	pm.configEndpoints()

	if c.Cfg.DeliveryType == "grpc" {
		pm.configGrpc()
	}

	return nil
}
