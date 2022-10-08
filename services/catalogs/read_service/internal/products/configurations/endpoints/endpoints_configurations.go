package endpoints

import (
	"context"
	"github.com/labstack/echo/v4"
	customEcho "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/custom_echo"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/delivery"
	gettingProductByIdV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/get_product_by_id/endpoints/v1"
	gettingProductsV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/getting_products/endpoints/v1"
	searchingProductsV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/searching_products/endpoints/v1"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/contracts"
)

func ConfigProductsEndpoints(ctx context.Context, routeBuilder *customEcho.RouteBuilder, infra contracts.InfrastructureConfigurations) {
	configV1Endpoints(routeBuilder, infra, ctx)
}

func configV1Endpoints(routeBuilder *customEcho.RouteBuilder, infra contracts.InfrastructureConfigurations, ctx context.Context) {
	routeBuilder.RegisterGroup("/api/v1", func(v1 *echo.Group) {
		group := v1.Group("/products")
		productEndpointBase := &delivery.ProductEndpointBase{
			ProductsGroup:                group,
			InfrastructureConfigurations: infra,
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
	})
}
