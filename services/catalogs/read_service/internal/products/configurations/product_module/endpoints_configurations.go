package product_module

import (
	"context"
	"github.com/labstack/echo/v4"
	customEcho "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/custom_echo"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/delivery"
	gettingProductByIdV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/get_product_by_id/endpoints/v1"
	gettingProductsV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/getting_products/endpoints/v1"
	searchingProductsV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/searching_products/endpoints/v1"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/configurations/infrastructure"
)

func (c *productsModuleConfigurator) configEndpoints(ctx context.Context) {
	configV1Endpoints(c.echoServer, c.InfrastructureConfigurations, ctx)
}

func configV1Endpoints(echoServer customEcho.EchoHttpServer, infra *infrastructure.InfrastructureConfigurations, ctx context.Context) {
	echoServer.ConfigGroup("/api/v1", func(v1 *echo.Group) {
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
