package configurations

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/delivery"
	gettingProductByIdV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/get_product_by_id/endpoints/v1"
	gettingProductsV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/getting_products/endpoints/v1"
	searchingProductsV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/searching_products/endpoints/v1"
)

func (c *productsModuleConfigurator) configEndpoints(ctx context.Context, group *echo.Group) {
	fmt.Print(c)

	productEndpointBase := &delivery.ProductEndpointBase{
		ProductsGroup:                group,
		InfrastructureConfigurations: c.InfrastructureConfigurations,
	}

	// GetProducts
	getProductsEndpoint := gettingProductsV1.NewGetProductsEndpoint(productEndpointBase)
	getProductsEndpoint.MapRoute()

	// SearchProducts
	searchProductsEndpoint := searchingProductsV1.NewSearchProductsEndpoint(productEndpointBase)
	searchProductsEndpoint.MapRoute()

	// GetProductById
	getProductByIdEndpoint := gettingProductByIdV1.NewGetProductByIdEndpoint(productEndpointBase)
	getProductByIdEndpoint.MapRoute()
}
