package mediatr

import (
	logger2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/producer"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing"
	"github.com/mehdihadeli/go-mediatr"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/contracts/data"
	createProductCommandV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/creating_product/v1/commands"
	createProductV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/creating_product/v1/dtos"
	deleteProductCommandV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/deleting_product/v1/commands"
	getProductByIdDtosV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/getting_product_by_id/v1/dtos"
	getProductByIdQueryV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/getting_product_by_id/v1/queries"
	getProductsDtosV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/getting_products/v1/dtos"
	getProductsQueryV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/getting_products/v1/queries"
	searchProductsDtosV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/searching_product/v1/dtos"
	searchProductsQueryV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/searching_product/v1/queries"
	updateProductCommandV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/updating_product/v1/commands"
)

func ConfigProductsMediator(
	logger logger2.Logger,
	uow data.CatalogUnitOfWork,
	productRepository data.ProductRepository,
	producer producer.Producer,
	tracer tracing.AppTracer,
) error {
	// https://stackoverflow.com/questions/72034479/how-to-implement-generic-interfaces
	err := mediatr.RegisterRequestHandler[*createProductCommandV1.CreateProduct, *createProductV1.CreateProductResponseDto](
		createProductCommandV1.NewCreateProductHandler(logger, uow, producer, tracer),
	)
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*getProductsQueryV1.GetProducts, *getProductsDtosV1.GetProductsResponseDto](
		getProductsQueryV1.NewGetProductsHandler(logger, productRepository, tracer),
	)
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*searchProductsQueryV1.SearchProducts, *searchProductsDtosV1.SearchProductsResponseDto](
		searchProductsQueryV1.NewSearchProductsHandler(logger, productRepository, tracer),
	)
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*updateProductCommandV1.UpdateProduct, *mediatr.Unit](
		updateProductCommandV1.NewUpdateProductHandler(logger, uow, producer, tracer),
	)
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*deleteProductCommandV1.DeleteProduct, *mediatr.Unit](
		deleteProductCommandV1.NewDeleteProductHandler(logger, uow, producer, tracer),
	)
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*getProductByIdQueryV1.GetProductById, *getProductByIdDtosV1.GetProductByIdResponseDto](
		getProductByIdQueryV1.NewGetProductByIdHandler(logger, productRepository, tracer),
	)
	if err != nil {
		return err
	}

	return nil
}
