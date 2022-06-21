package products

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mediatr"
	v1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/features/creating_product/v1"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/infrastructure/repositories"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/shared"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/shared/server/configurations"
)

type ProductsModuleConfigurator interface {
	ConfigureProductsModule() error
}

type productsModuleConfigurator struct {
	server *configurations.Server
}

func NewProductsModuleConfigurator(server *configurations.Server) *productsModuleConfigurator {
	return &productsModuleConfigurator{server: server}
}

func (p *productsModuleConfigurator) ConfigureProductsModule() error {

	productRepo := repositories.NewPostgresProductRepository(p.server.Log, p.server.Cfg, p.server.PgConn)
	m, err := shared.NewMediator(p.server.Log, p.server.Cfg, productRepo, p.server.KafkaProducer)

	if err != nil {
		return err
	}

	p.configEndpoints(m)

	return nil
}

func (p *productsModuleConfigurator) configEndpoints(m *mediatr.Mediator) {

	v1.NewCreteProductEndpoint(p.server.Echo, p.server.Log, p.server.Cfg, m, p.server.Validator, p.server.Metrics)
}
