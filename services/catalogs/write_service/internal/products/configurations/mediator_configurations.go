package configurations

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mediatr"
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
	err := mediatr.RegisterHandler[*creatingProductV1.CreateProduct, *creatingProductsDtos.CreateProductResponseDto](creatingProductV1.NewCreateProductHandler(c.Log, c.Cfg, pgRepo, c.KafkaProducer))
	if err != nil {
		return err
	}

	err = mediatr.RegisterHandler[*gettingProductsV1.GetProducts, *gettingProductsDtos.GetProductsResponseDto](gettingProductsV1.NewGetProductsHandler(c.Log, c.Cfg, pgRepo))
	if err != nil {
		return err
	}

	err = mediatr.RegisterHandler[*searchingProductsV1.SearchProducts, *searchingProductsDtos.SearchProductsResponseDto](searchingProductsV1.NewSearchProductsHandler(c.Log, c.Cfg, pgRepo))
	if err != nil {
		return err
	}

	err = mediatr.RegisterHandler[*updatingProductV1.UpdateProduct, *mediatr.Unit](updatingProductV1.NewUpdateProductHandler(c.Log, c.Cfg, pgRepo, c.KafkaProducer))
	if err != nil {
		return err
	}

	err = mediatr.RegisterHandler[*deletingProductV1.DeleteProduct, *mediatr.Unit](deletingProductV1.NewDeleteProductHandler(c.Log, c.Cfg, pgRepo, c.KafkaProducer))
	if err != nil {
		return err
	}

	err = mediatr.RegisterHandler[*geettingProductByIdV1.GetProductById, *gettingProductByIdDtos.GetProductByIdResponseDto](geettingProductByIdV1.NewGetProductByIdHandler(c.Log, c.Cfg, pgRepo))
	if err != nil {
		return err
	}

	return nil
}
