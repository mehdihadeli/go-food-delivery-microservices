package configurations

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/creating_product"
	creating_products_dtos "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/creating_product/dtos"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/getting_products"
	getting_products_dtos "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/getting_products/dtos"
	"github.com/pkg/errors"
)

func (c *productsModuleConfigurator) configProductsMediator(pgRepo contracts.ProductRepository) error {

	err := mediatr.RegisterHandler[*creating_product.CreateProduct, *creating_products_dtos.CreateProductResponseDto](creating_product.NewCreateProductHandler(c.Log, c.Cfg, pgRepo, c.KafkaProducer))
	if err != nil {
		return errors.Wrap(err, "error while registering handlers in the mediator")
	}

	err = mediatr.RegisterHandler[*getting_products.GetProducts, *getting_products_dtos.GetProductsResponseDto](getting_products.NewGetProductsHandler(c.Log, c.Cfg, pgRepo))
	if err != nil {
		return errors.Wrap(err, "error while registering handlers in the mediator")
	}

	return nil
}
