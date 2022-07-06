package configurations

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/creating_product"
	"github.com/pkg/errors"
)

func (c *productsModuleConfigurator) configProductsMediator(pgRepo contracts.ProductRepository) (*mediatr.Mediator, error) {

	md := mediatr.New()

	err := md.Register(creating_product.NewCreateProductHandler(c.Log, c.Cfg, pgRepo, c.KafkaProducer))

	if err != nil {
		return nil, errors.Wrap(err, "error while registering handlers in the mediator")
	}

	return &md, nil
}
