package configurations

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/creating_product"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/deleting_products"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/getting_products"
	getting_products_dtos "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/getting_products/dtos"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/updating_products"
	"github.com/pkg/errors"
)

func (c *productsModuleConfigurator) configProductsMediator(repository contracts.ProductRepository) error {

	err := mediatr.RegisterHandler[*creating_product.CreateProduct, *creating_product.CreateProductResponseDto](creating_product.NewCreateProductHandler(c.Log, c.Cfg, repository))
	if err != nil {
		return errors.Wrap(err, "error while registering handlers in the mediator")
	}

	err = mediatr.RegisterHandler[*deleting_products.DeleteProduct, *mediatr.Unit](deleting_products.NewDeleteProductHandler(c.Log, c.Cfg, repository))
	if err != nil {
		return errors.Wrap(err, "error while registering handlers in the mediator")
	}

	err = mediatr.RegisterHandler[*updating_products.UpdateProduct, *mediatr.Unit](updating_products.NewUpdateProductHandler(c.Log, c.Cfg, repository))
	if err != nil {
		return errors.Wrap(err, "error while registering handlers in the mediator")
	}

	err = mediatr.RegisterHandler[*getting_products.GetProducts, *getting_products_dtos.GetProductsResponseDto](getting_products.NewGetProductsHandler(c.Log, c.Cfg, repository))
	if err != nil {
		return errors.Wrap(err, "error while registering handlers in the mediator")
	}

	return nil
}
