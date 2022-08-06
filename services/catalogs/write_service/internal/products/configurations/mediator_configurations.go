package configurations

import (
	"github.com/mehdihadeli/go-mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts"
	creatingProductV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/creating_product/commands/v1"
	creatingProductsDtos "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/creating_product/dtos"
	deletingProductV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/deleting_product/commands/v1"
	gettingProductByIdDtos "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_product_by_id/dtos"
	geettingProductByIdV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_product_by_id/queries/v1"
	gettingProductsDtos "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_products/dtos"
	gettingProductsV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_products/queries/v1"
	searchingProductsDtos "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/searching_product/dtos"
	searchingProductsV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/searching_product/queries/v1"
	updatingProductV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/updating_product/commands/v1"
)

func (c *productsModuleConfigurator) configProductsMediator(pgRepo contracts.ProductRepository) error {

	//https://stackoverflow.com/questions/72034479/how-to-implement-generic-interfaces
	err := mediatr.RegisterRequestHandler[*creatingProductV1.CreateProductCommand, *creatingProductsDtos.CreateProductResponseDto](creatingProductV1.NewCreateProductCommandHandler(c.Log, c.Cfg, pgRepo, c.KafkaProducer))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*gettingProductsV1.GetProductsQuery, *gettingProductsDtos.GetProductsResponseDto](gettingProductsV1.NewGetProductsQueryHandler(c.Log, c.Cfg, pgRepo))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*searchingProductsV1.SearchProductsQuery, *searchingProductsDtos.SearchProductsResponseDto](searchingProductsV1.NewSearchProductsQueryHandler(c.Log, c.Cfg, pgRepo))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*updatingProductV1.UpdateProductCommand, *mediatr.Unit](updatingProductV1.NewUpdateProductCommandHandler(c.Log, c.Cfg, pgRepo, c.KafkaProducer))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*deletingProductV1.DeleteProductCommand, *mediatr.Unit](deletingProductV1.NewDeleteProductCommandHandler(c.Log, c.Cfg, pgRepo, c.KafkaProducer))
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*geettingProductByIdV1.GetProductByIdQuery, *gettingProductByIdDtos.GetProductByIdResponseDto](geettingProductByIdV1.NewGetProductByIdQueryHandler(c.Log, c.Cfg, pgRepo))
	if err != nil {
		return err
	}

	return nil
}
