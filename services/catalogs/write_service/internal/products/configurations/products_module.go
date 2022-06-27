package configurations

import (
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mediatr"
	product_service "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts/grpc/service_clients"
	repositories_contract "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts/repositories"
	repositories_imp "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/data/repositories"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/delivery/grpc"
	creating_product "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/creating_product/endpoints/v1"
	deleting_product "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/deleting_product/endpoints/v1"
	gettting_product_by_id "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_product_by_id/endpoints/v1"
	getting_products "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_products/endpoints/v1"
	searching_products "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/searching_product/endpoints/v1"
	updating_product "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/updating_product/endpoints/v1"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/configurations"
)

type ProductsModuleConfigurator interface {
	ConfigureProductsModule() error
}

type productsModuleConfigurator struct {
	infrastructure *configurations.Infrastructure
}

type ProductModule struct {
	Infrastructure    *configurations.Infrastructure
	Mediator          *mediatr.Mediator
	ProductRepository repositories_contract.ProductRepository
	ProductsGroup     *echo.Group
}

func NewProductsModuleConfigurator(infrastructure *configurations.Infrastructure) *productsModuleConfigurator {
	return &productsModuleConfigurator{infrastructure: infrastructure}
}

func (c *productsModuleConfigurator) ConfigureProductsModule() error {

	pm := ProductModule{Infrastructure: c.infrastructure}

	v1 := c.infrastructure.Echo.Group("/api/v1")
	pm.ProductsGroup = v1.Group("/" + c.infrastructure.Cfg.Http.ProductsPath)

	pm.ProductRepository = repositories_imp.NewPostgresProductRepository(c.infrastructure.Log, c.infrastructure.Cfg, c.infrastructure.PgConn, c.infrastructure.Gorm)
	m, err := shared.NewCatalogsMediator(c.infrastructure.Log, c.infrastructure.Cfg, pm.ProductRepository, c.infrastructure.KafkaProducer)

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

	// CreateNewProduct
	createProductEndpoint := creating_product.NewCreteProductEndpoint(pm.Infrastructure, pm.Mediator, pm.ProductsGroup, pm.ProductRepository)
	createProductEndpoint.MapRoute()

	// UpdateProduct
	updateProductEndpoint := updating_product.NewUpdateProductEndpoint(pm.Infrastructure, pm.Mediator, pm.ProductsGroup, pm.ProductRepository)
	updateProductEndpoint.MapRoute()

	// GetProducts
	getProductsEndpoint := getting_products.NewGetProductsEndpoint(pm.Infrastructure, pm.Mediator, pm.ProductsGroup, pm.ProductRepository)
	getProductsEndpoint.MapRoute()

	// SearchProducts
	searchProducts := searching_products.NewSearchProductsEndpoint(pm.Infrastructure, pm.Mediator, pm.ProductsGroup, pm.ProductRepository)
	searchProducts.MapRoute()

	// GetProductById
	getProductByIdEndpoint := gettting_product_by_id.NewGetProductByIdEndpoint(pm.Infrastructure, pm.Mediator, pm.ProductsGroup, pm.ProductRepository)
	getProductByIdEndpoint.MapRoute()

	// DeleteProduct
	deleteProductEndpoint := deleting_product.NewDeleteProductEndpoint(pm.Infrastructure, pm.Mediator, pm.ProductsGroup, pm.ProductRepository)
	deleteProductEndpoint.MapRoute()
}

func (pm *ProductModule) configGrpc() {
	productGrpcService := grpc.NewProductGrpcService(pm.Infrastructure, pm.Mediator, pm.ProductRepository)
	product_service.RegisterProductsServiceServer(pm.Infrastructure.GrpcServer, productGrpcService)
}
