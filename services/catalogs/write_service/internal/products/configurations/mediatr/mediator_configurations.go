package mediatr

import (
	"github.com/mehdihadeli/go-mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/producer"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts/data"
	createProductCommandV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/creating_product/v1/commands"
	createProductsDtosV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/creating_product/v1/dtos"
	deleteProductCommandV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/deleting_product/v1/commands"
	getProductByIdDtosV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_product_by_id/v1/dtos"
	getProductByIdQueryV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_product_by_id/v1/queries"
	getProductsDtosV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_products/v1/dtos"
	getProductsQueryV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_products/v1/queries"
	searchProductsDtosV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/searching_product/v1/dtos"
	searchProductsQueryV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/searching_product/v1/queries"
	updateProductCommandV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/updating_product/v1/commands"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/contracts"
)

func ConfigProductsMediator(uow data.CatalogUnitOfWork, productRepository data.ProductRepository, infra *contracts.InfrastructureConfigurations, producer producer.Producer) error {
	// https://stackoverflow.com/questions/72034479/how-to-implement-generic-interfaces
	err := mediatr.RegisterRequestHandler[*createProductCommandV1.CreateProduct, *createProductsDtosV1.CreateProductResponseDto](createProductCommandV1.NewCreateProductHandler(infra.Log, infra.Cfg, uow, producer))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*getProductsQueryV1.GetProducts, *getProductsDtosV1.GetProductsResponseDto](getProductsQueryV1.NewGetProductsHandler(infra.Log, infra.Cfg, productRepository))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*searchProductsQueryV1.SearchProducts, *searchProductsDtosV1.SearchProductsResponseDto](searchProductsQueryV1.NewSearchProductsHandler(infra.Log, infra.Cfg, productRepository))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*updateProductCommandV1.UpdateProduct, *mediatr.Unit](updateProductCommandV1.NewUpdateProductHandler(infra.Log, infra.Cfg, uow, producer))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*deleteProductCommandV1.DeleteProduct, *mediatr.Unit](deleteProductCommandV1.NewDeleteProductHandler(infra.Log, infra.Cfg, uow, producer))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*getProductByIdQueryV1.GetProductById, *getProductByIdDtosV1.GetProductByIdResponseDto](getProductByIdQueryV1.NewGetProductByIdHandler(infra.Log, infra.Cfg, productRepository))
	if err != nil {
		return err
	}

	return nil
}
