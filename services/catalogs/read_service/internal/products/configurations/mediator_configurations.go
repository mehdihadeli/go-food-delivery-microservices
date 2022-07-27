package configurations

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/creating_product"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/creating_product/commands/v1"
	deletingProductV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/deleting_products/commands/v1"
	gettingProductByIdDtos "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/get_product_by_id/dtos"
	gettingProductByIdV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/get_product_by_id/queries/v1"
	gettingProductsDtos "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/getting_products/dtos"
	gettingProductsV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/getting_products/queries/v1"
	searchingProductsDtos "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/searching_products/dtos"
	searchingProductV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/searching_products/queries/v1"
	updatingProductV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/updating_products/commands/v1"
	"github.com/pkg/errors"
)

func (c *productsModuleConfigurator) configProductsMediator(mongoRepository contracts.ProductRepository, redisRepository contracts.ProductCacheRepository) error {

	err := mediatr.RegisterHandler[*v1.CreateProduct, *creating_product.CreateProductResponseDto](v1.NewCreateProductHandler(c.Log, c.Cfg, mongoRepository, redisRepository))
	if err != nil {
		return errors.Wrap(err, "error while registering handlers in the mediator")
	}

	err = mediatr.RegisterHandler[*deletingProductV1.DeleteProduct, *mediatr.Unit](deletingProductV1.NewDeleteProductHandler(c.Log, c.Cfg, mongoRepository, redisRepository))
	if err != nil {
		return errors.Wrap(err, "error while registering handlers in the mediator")
	}

	err = mediatr.RegisterHandler[*updatingProductV1.UpdateProduct, *mediatr.Unit](updatingProductV1.NewUpdateProductHandler(c.Log, c.Cfg, mongoRepository, redisRepository))
	if err != nil {
		return errors.Wrap(err, "error while registering handlers in the mediator")
	}

	err = mediatr.RegisterHandler[*gettingProductsV1.GetProducts, *gettingProductsDtos.GetProductsResponseDto](gettingProductsV1.NewGetProductsHandler(c.Log, c.Cfg, mongoRepository))
	if err != nil {
		return errors.Wrap(err, "error while registering handlers in the mediator")
	}

	err = mediatr.RegisterHandler[*searchingProductV1.SearchProducts, *searchingProductsDtos.SearchProductsResponseDto](searchingProductV1.NewSearchProductsHandler(c.Log, c.Cfg, mongoRepository))
	if err != nil {
		return errors.Wrap(err, "error while registering handlers in the mediator")
	}

	err = mediatr.RegisterHandler[*gettingProductByIdV1.GetProductById, *gettingProductByIdDtos.GetProductByIdResponseDto](gettingProductByIdV1.NewGetProductByIdHandler(c.Log, c.Cfg, mongoRepository, redisRepository))
	if err != nil {
		return errors.Wrap(err, "error while registering handlers in the mediator")
	}

	return nil
}
