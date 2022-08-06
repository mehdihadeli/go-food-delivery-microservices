package configurations

import (
	"github.com/mehdihadeli/go-mediatr"
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

	err := mediatr.RegisterRequestHandler[*v1.CreateProductCommand, *creating_product.CreateProductResponseDto](v1.NewCreateProductCommandHandler(c.Log, c.Cfg, mongoRepository, redisRepository))
	if err != nil {
		return errors.Wrap(err, "error while registering handlers in the mediator")
	}

	err = mediatr.RegisterRequestHandler[*deletingProductV1.DeleteProductCommand, *mediatr.Unit](deletingProductV1.NewDeleteProductCommandHandler(c.Log, c.Cfg, mongoRepository, redisRepository))
	if err != nil {
		return errors.Wrap(err, "error while registering handlers in the mediator")
	}

	err = mediatr.RegisterRequestHandler[*updatingProductV1.UpdateProductCommand, *mediatr.Unit](updatingProductV1.NewUpdateProductCommandHandler(c.Log, c.Cfg, mongoRepository, redisRepository))
	if err != nil {
		return errors.Wrap(err, "error while registering handlers in the mediator")
	}

	err = mediatr.RegisterRequestHandler[*gettingProductsV1.GetProductsQuery, *gettingProductsDtos.GetProductsResponseDto](gettingProductsV1.NewGetProductsQueryHandler(c.Log, c.Cfg, mongoRepository))
	if err != nil {
		return errors.Wrap(err, "error while registering handlers in the mediator")
	}

	err = mediatr.RegisterRequestHandler[*searchingProductV1.SearchProductsQuery, *searchingProductsDtos.SearchProductsResponseDto](searchingProductV1.NewSearchProductsQueryHandler(c.Log, c.Cfg, mongoRepository))
	if err != nil {
		return errors.Wrap(err, "error while registering handlers in the mediator")
	}

	err = mediatr.RegisterRequestHandler[*gettingProductByIdV1.GetProductByIdQuery, *gettingProductByIdDtos.GetProductByIdResponseDto](gettingProductByIdV1.NewGetProductByIdQueryHandler(c.Log, c.Cfg, mongoRepository, redisRepository))
	if err != nil {
		return errors.Wrap(err, "error while registering handlers in the mediator")
	}

	return nil
}
