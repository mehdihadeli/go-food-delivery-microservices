package endpoints

import (
	"context"
	"github.com/labstack/echo/v4"
	customEcho "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/custom_echo"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/bus"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/delivery"
	gettingProductByIdV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/get_product_by_id/endpoints/v1"
	gettingProductsV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/getting_products/endpoints/v1"
	searchingProductsV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/searching_products/endpoints/v1"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/contracts"
)

func ConfigProductsEndpoints(ctx context.Context, routeBuilder *customEcho.RouteBuilder, infra *contracts.InfrastructureConfigurations, bus bus.Bus, metrics *contracts.CatalogsMetrics) {
	configV1Endpoints(ctx, routeBuilder, infra, bus, metrics)
}

func configV1Endpoints(ctx context.Context, routeBuilder *customEcho.RouteBuilder, infra *contracts.InfrastructureConfigurations, bus bus.Bus, metrics *contracts.CatalogsMetrics) {
	routeBuilder.RegisterGroup("/api/v1", func(v1 *echo.Group) {
		group := v1.Group("/products")
		productEndpointBase := delivery.NewProductEndpointBase(infra, group, bus, metrics)

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
