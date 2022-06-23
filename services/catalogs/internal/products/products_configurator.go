package products

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mediatr"
	v1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/features/creating_product/v1"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/infrastructure/repositories"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/shared"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/shared/configurations"
)

type ProductsModuleConfigurator interface {
	ConfigureProductsModule() error
}

type productsModuleConfigurator struct {
	infrastructure *configurations.Infrastructure
}

func NewProductsModuleConfigurator(infrastructure *configurations.Infrastructure) *productsModuleConfigurator {
	return &productsModuleConfigurator{infrastructure: infrastructure}
}

func (p *productsModuleConfigurator) ConfigureProductsModule() error {

	productRepo := repositories.NewPostgresProductRepository(p.infrastructure.Log, p.infrastructure.Cfg, p.infrastructure.PgConn)
	m, err := shared.NewMediator(p.infrastructure.Log, p.infrastructure.Cfg, productRepo, p.infrastructure.KafkaProducer)

	if err != nil {
		return err
	}

	p.configEndpoints(m)

	return nil
}

func (p *productsModuleConfigurator) configEndpoints(m *mediatr.Mediator) {

	v1.NewCreteProductEndpoint(p.infrastructure.Echo, p.infrastructure.Log, p.infrastructure.Cfg, m, p.infrastructure.Validator, p.infrastructure.Metrics)
}
