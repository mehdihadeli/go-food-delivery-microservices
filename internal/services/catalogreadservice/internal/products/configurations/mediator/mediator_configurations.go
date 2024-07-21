package mediator

import (
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/otel/tracing"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogreadservice/internal/products/contracts/data"
	v1 "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogreadservice/internal/products/features/creating_product/v1"
	createProductDtosV1 "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogreadservice/internal/products/features/creating_product/v1/dtos"
	deleteProductCommandV1 "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogreadservice/internal/products/features/deleting_products/v1/commands"
	getProductByIdDtosV1 "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogreadservice/internal/products/features/get_product_by_id/v1/dtos"
	getProductByIdQueryV1 "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogreadservice/internal/products/features/get_product_by_id/v1/queries"
	getProductsDtoV1 "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogreadservice/internal/products/features/getting_products/v1/dtos"
	getProductsQueryV1 "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogreadservice/internal/products/features/getting_products/v1/queries"
	searchProductsDtosV1 "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogreadservice/internal/products/features/searching_products/v1/dtos"
	searchProductsQueryV1 "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogreadservice/internal/products/features/searching_products/v1/queries"
	updateProductCommandV1 "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogreadservice/internal/products/features/updating_products/v1/commands"

	"emperror.dev/errors"
	"github.com/mehdihadeli/go-mediatr"
)

func ConfigProductsMediator(
	logger logger.Logger,
	mongoProductRepository data.ProductRepository,
	cacheProductRepository data.ProductCacheRepository,
	tracer tracing.AppTracer,
) error {
	err := mediatr.RegisterRequestHandler[*v1.CreateProduct, *createProductDtosV1.CreateProductResponseDto](
		v1.NewCreateProductHandler(
			logger,
			mongoProductRepository,
			cacheProductRepository,
			tracer,
		),
	)
	if err != nil {
		return errors.WrapIf(err, "error while registering handlers in the mediator")
	}

	err = mediatr.RegisterRequestHandler[*deleteProductCommandV1.DeleteProduct, *mediatr.Unit](
		deleteProductCommandV1.NewDeleteProductHandler(
			logger,
			mongoProductRepository,
			cacheProductRepository,
			tracer,
		),
	)
	if err != nil {
		return errors.WrapIf(err, "error while registering handlers in the mediator")
	}

	err = mediatr.RegisterRequestHandler[*updateProductCommandV1.UpdateProduct, *mediatr.Unit](
		updateProductCommandV1.NewUpdateProductHandler(
			logger,
			mongoProductRepository,
			cacheProductRepository,
			tracer,
		),
	)
	if err != nil {
		return errors.WrapIf(err, "error while registering handlers in the mediator")
	}

	err = mediatr.RegisterRequestHandler[*getProductsQueryV1.GetProducts, *getProductsDtoV1.GetProductsResponseDto](
		getProductsQueryV1.NewGetProductsHandler(logger, mongoProductRepository, tracer),
	)
	if err != nil {
		return errors.WrapIf(err, "error while registering handlers in the mediator")
	}

	err = mediatr.RegisterRequestHandler[*searchProductsQueryV1.SearchProducts, *searchProductsDtosV1.SearchProductsResponseDto](
		searchProductsQueryV1.NewSearchProductsHandler(
			logger,
			mongoProductRepository,
			tracer,
		),
	)
	if err != nil {
		return errors.WrapIf(err, "error while registering handlers in the mediator")
	}

	err = mediatr.RegisterRequestHandler[*getProductByIdQueryV1.GetProductById, *getProductByIdDtosV1.GetProductByIdResponseDto](
		getProductByIdQueryV1.NewGetProductByIdHandler(
			logger,
			mongoProductRepository,
			cacheProductRepository,
			tracer,
		),
	)
	if err != nil {
		return errors.WrapIf(err, "error while registering handlers in the mediator")
	}

	return nil
}
