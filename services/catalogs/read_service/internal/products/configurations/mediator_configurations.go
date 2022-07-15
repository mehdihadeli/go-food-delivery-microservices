package configurations

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/creating_product"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/deleting_products"
	getting_product_by_id "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/get_product_by_id"
	getting_product_by_id_dtos "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/get_product_by_id/dtos"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/getting_products"
	getting_products_dtos "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/getting_products/dtos"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/searching_products"
	searching_products_dtos "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/searching_products/dtos"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/updating_products"
	"github.com/pkg/errors"
)

func (c *productsModuleConfigurator) configProductsMediator(mongoRepository contracts.ProductRepository, redisRepository contracts.ProductCacheRepository) error {

	err := mediatr.RegisterHandler[*creating_product.CreateProduct, *creating_product.CreateProductResponseDto](creating_product.NewCreateProductHandler(c.Log, c.Cfg, mongoRepository, redisRepository))
	if err != nil {
		return errors.Wrap(err, "error while registering handlers in the mediator")
	}

	err = mediatr.RegisterHandler[*deleting_products.DeleteProduct, *mediatr.Unit](deleting_products.NewDeleteProductHandler(c.Log, c.Cfg, mongoRepository, redisRepository))
	if err != nil {
		return errors.Wrap(err, "error while registering handlers in the mediator")
	}

	err = mediatr.RegisterHandler[*updating_products.UpdateProduct, *mediatr.Unit](updating_products.NewUpdateProductHandler(c.Log, c.Cfg, mongoRepository, redisRepository))
	if err != nil {
		return errors.Wrap(err, "error while registering handlers in the mediator")
	}

	err = mediatr.RegisterHandler[*getting_products.GetProducts, *getting_products_dtos.GetProductsResponseDto](getting_products.NewGetProductsHandler(c.Log, c.Cfg, mongoRepository))
	if err != nil {
		return errors.Wrap(err, "error while registering handlers in the mediator")
	}

	err = mediatr.RegisterHandler[*searching_products.SearchProducts, *searching_products_dtos.SearchProductsResponseDto](searching_products.NewSearchProductsHandler(c.Log, c.Cfg, mongoRepository))
	if err != nil {
		return errors.Wrap(err, "error while registering handlers in the mediator")
	}

	err = mediatr.RegisterHandler[*getting_product_by_id.GetProductById, *getting_product_by_id_dtos.GetProductByIdResponseDto](getting_product_by_id.NewGetProductByIdHandler(c.Log, c.Cfg, mongoRepository, redisRepository))
	if err != nil {
		return errors.Wrap(err, "error while registering handlers in the mediator")
	}

	return nil
}
