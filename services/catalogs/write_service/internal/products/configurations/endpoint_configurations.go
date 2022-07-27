package configurations

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/delivery"
	creatingProductV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/creating_product/endpoints/v1"
	deletingProductV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/deleting_product/endpoints/v1"
	gettingProductByIdV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_product_by_id/endpoints/v1"
	gettingProductsV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_products/endpoints/v1"
	searchingProductsV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/searching_product/endpoints/v1"
	updatingProductV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/updating_product/endpoints/v1"
)

func (c *productsModuleConfigurator) configEndpoints(ctx context.Context, group *echo.Group) {

	productEndpointBase := &delivery.ProductEndpointBase{
		ProductsGroup:               group,
		InfrastructureConfiguration: c.InfrastructureConfiguration,
	}
	// CreateNewProduct
	createProductEndpoint := creatingProductV1.NewCreteProductEndpoint(productEndpointBase)
	createProductEndpoint.MapRoute()

	// UpdateProduct
	updateProductEndpoint := updatingProductV1.NewUpdateProductEndpoint(productEndpointBase)
	updateProductEndpoint.MapRoute()

	// GetProducts
	getProductsEndpoint := gettingProductsV1.NewGetProductsEndpoint(productEndpointBase)
	getProductsEndpoint.MapRoute()

	// SearchProducts
	searchProducts := searchingProductsV1.NewSearchProductsEndpoint(productEndpointBase)
	searchProducts.MapRoute()

	// GetProductById
	getProductByIdEndpoint := gettingProductByIdV1.NewGetProductByIdEndpoint(productEndpointBase)
	getProductByIdEndpoint.MapRoute()

	// DeleteProduct
	deleteProductEndpoint := deletingProductV1.NewDeleteProductEndpoint(productEndpointBase)
	deleteProductEndpoint.MapRoute()
}
