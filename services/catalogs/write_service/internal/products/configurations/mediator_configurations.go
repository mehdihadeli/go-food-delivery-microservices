package configurations

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/creating_product"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/deleting_product"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_product_by_id"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_products"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/searching_product"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/updating_product"
	"github.com/pkg/errors"
)

func (c *productsModuleConfigurator) configProductsMediator(pgRepo contracts.ProductRepository) (*mediatr.Mediator, error) {

	md := mediatr.New()

	err := md.Register(
		creating_product.NewCreateProductHandler(c.Log, c.Cfg, pgRepo, c.KafkaProducer),
		updating_product.NewUpdateProductHandler(c.Log, c.Cfg, pgRepo, c.KafkaProducer),
		deleting_product.NewDeleteProductHandler(c.Log, c.Cfg, pgRepo, c.KafkaProducer),
		getting_product_by_id.NewGetProductByIdHandler(c.Log, c.Cfg, pgRepo),
		getting_products.NewGetProductsHandler(c.Log, c.Cfg, pgRepo),
		searching_product.NewSearchProductsHandler(c.Log, c.Cfg, pgRepo),
	)

	if err != nil {
		return nil, errors.Wrap(err, "error while registering handlers in the mediator")
	}

	return &md, nil
}
