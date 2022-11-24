package mediator

import (
	"emperror.dev/errors"
	"github.com/mehdihadeli/go-mediatr"

	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/bus"
	contracts2 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/contracts"
	createProductCommandV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/creating_product/v1/commands"
	createProductDtosV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/creating_product/v1/dtos"
	deleteProductCommandV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/deleting_products/v1/commands"
	getProductByIdDtosV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/get_product_by_id/v1/dtos"
	getProductByIdQueryV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/get_product_by_id/v1/queries"
	getProductsDtoV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/getting_products/v1/dtos"
	getProductsQueryV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/getting_products/v1/queries"
	searchProductsDtosV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/searching_products/v1/dtos"
	searchProductsQueryV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/searching_products/v1/queries"
	updateProductCommandV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/updating_products/v1/commands"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/contracts"
)

func ConfigProductsMediator(infra *contracts.InfrastructureConfigurations, mongoProductRepository contracts2.ProductRepository, cacheProductRepository contracts2.ProductCacheRepository, bus bus.Bus) error {
	err := mediatr.RegisterRequestHandler[*createProductCommandV1.CreateProduct, *createProductDtosV1.CreateProductResponseDto](createProductCommandV1.NewCreateProductHandler(infra.Log, infra.Cfg, mongoProductRepository, cacheProductRepository))
	if err != nil {
		return errors.WrapIf(err, "error while registering handlers in the mediator")
	}

	err = mediatr.RegisterRequestHandler[*deleteProductCommandV1.DeleteProduct, *mediatr.Unit](deleteProductCommandV1.NewDeleteProductHandler(infra.Log, infra.Cfg, mongoProductRepository, cacheProductRepository))
	if err != nil {
		return errors.WrapIf(err, "error while registering handlers in the mediator")
	}

	err = mediatr.RegisterRequestHandler[*updateProductCommandV1.UpdateProduct, *mediatr.Unit](updateProductCommandV1.NewUpdateProductHandler(infra.Log, infra.Cfg, mongoProductRepository, cacheProductRepository))
	if err != nil {
		return errors.WrapIf(err, "error while registering handlers in the mediator")
	}

	err = mediatr.RegisterRequestHandler[*getProductsQueryV1.GetProducts, *getProductsDtoV1.GetProductsResponseDto](getProductsQueryV1.NewGetProductsHandler(infra.Log, infra.Cfg, mongoProductRepository))
	if err != nil {
		return errors.WrapIf(err, "error while registering handlers in the mediator")
	}

	err = mediatr.RegisterRequestHandler[*searchProductsQueryV1.SearchProducts, *searchProductsDtosV1.SearchProductsResponseDto](searchProductsQueryV1.NewSearchProductsHandler(infra.Log, infra.Cfg, mongoProductRepository))
	if err != nil {
		return errors.WrapIf(err, "error while registering handlers in the mediator")
	}

	err = mediatr.RegisterRequestHandler[*getProductByIdQueryV1.GetProductById, *getProductByIdDtosV1.GetProductByIdResponseDto](getProductByIdQueryV1.NewGetProductByIdHandler(infra.Log, infra.Cfg, mongoProductRepository, cacheProductRepository))
	if err != nil {
		return errors.WrapIf(err, "error while registering handlers in the mediator")
	}

	return nil
}
