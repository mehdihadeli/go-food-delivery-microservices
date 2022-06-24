package configurations

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mediatr"
	product_service "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/contracts/grpc/service_clients"
	repositories_contract "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/contracts/repositories"
	repositories_imp "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/data/repositories"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/delivery/grpc"
	v1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/features/creating_product/endpoints/v1"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/shared"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/shared/configurations"
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
}

func NewProductsModuleConfigurator(infrastructure *configurations.Infrastructure) *productsModuleConfigurator {
	return &productsModuleConfigurator{infrastructure: infrastructure}
}

func (c *productsModuleConfigurator) ConfigureProductsModule() error {

	pm := ProductModule{Infrastructure: c.infrastructure}

	pm.ProductRepository = repositories_imp.NewPostgresProductRepository(c.infrastructure.Log, c.infrastructure.Cfg, c.infrastructure.PgConn)
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
	createProductEndpoint := v1.NewCreteProductEndpoint(pm.Infrastructure, pm.Mediator, pm.ProductRepository)
	createProductEndpoint.MapRoute()

	// UpdateProduct
}

func (pm *ProductModule) configGrpc() {

	productGrpcService := grpc.NewProductGrpcService(pm.Infrastructure, pm.Mediator, pm.ProductRepository)
	product_service.RegisterProductsServiceServer(pm.Infrastructure.GrpcServer, productGrpcService)
}
