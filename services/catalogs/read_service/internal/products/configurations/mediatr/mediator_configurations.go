package mediatr

import (
	"emperror.dev/errors"
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
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/configurations/infrastructure"
)

func ConfigProductsMediator(mongoRepository contracts.ProductRepository, redisRepository contracts.ProductCacheRepository, infra *infrastructure.InfrastructureConfigurations) error {
	err := mediatr.RegisterRequestHandler[*v1.CreateProduct, *creating_product.CreateProductResponseDto](v1.NewCreateProductHandler(infra.Log, infra.Cfg, mongoRepository, redisRepository))
	if err != nil {
		return errors.WrapIf(err, "error while registering handlers in the mediator")
	}

	err = mediatr.RegisterRequestHandler[*deletingProductV1.DeleteProduct, *mediatr.Unit](deletingProductV1.NewDeleteProductCommand(infra.Log, infra.Cfg, mongoRepository, redisRepository))
	if err != nil {
		return errors.WrapIf(err, "error while registering handlers in the mediator")
	}

	err = mediatr.RegisterRequestHandler[*updatingProductV1.UpdateProduct, *mediatr.Unit](updatingProductV1.NewUpdateProductHandler(infra.Log, infra.Cfg, mongoRepository, redisRepository))
	if err != nil {
		return errors.WrapIf(err, "error while registering handlers in the mediator")
	}

	err = mediatr.RegisterRequestHandler[*gettingProductsV1.GetProducts, *gettingProductsDtos.GetProductsResponseDto](gettingProductsV1.NewGetProductsHandler(infra.Log, infra.Cfg, mongoRepository))
	if err != nil {
		return errors.WrapIf(err, "error while registering handlers in the mediator")
	}

	err = mediatr.RegisterRequestHandler[*searchingProductV1.SearchProducts, *searchingProductsDtos.SearchProductsResponseDto](searchingProductV1.NewSearchProductsHandler(infra.Log, infra.Cfg, mongoRepository))
	if err != nil {
		return errors.WrapIf(err, "error while registering handlers in the mediator")
	}

	err = mediatr.RegisterRequestHandler[*gettingProductByIdV1.GetProductById, *gettingProductByIdDtos.GetProductByIdResponseDto](gettingProductByIdV1.NewGetProductByIdHandler(infra.Log, infra.Cfg, mongoRepository, redisRepository))
	if err != nil {
		return errors.WrapIf(err, "error while registering handlers in the mediator")
	}

	return nil
}
