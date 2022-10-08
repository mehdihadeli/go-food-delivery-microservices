package endpoints

import (
	"context"
	"github.com/labstack/echo/v4"
	customEcho "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/custom_echo"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/bus"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/delivery"
	creatingProductV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/creating_product/endpoints/v1"
	deletingProductV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/deleting_product/endpoints/v1"
	gettingProductByIdV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_product_by_id/endpoints/v1"
	gettingProductsV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_products/endpoints/v1"
	searchingProductV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/searching_product/endpoints/v1"
	updatingProductV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/updating_product/endpoints/v1"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/contracts"
)

func ConfigProductsEndpoints(ctx context.Context, routeBuilder *customEcho.RouteBuilder, infra contracts.InfrastructureConfigurations, bus bus.Bus, metrics contracts.CatalogsMetrics) {
	configV1Endpoints(ctx, routeBuilder, infra, bus, metrics)
}

func configV1Endpoints(ctx context.Context, routeBuilder *customEcho.RouteBuilder, infra contracts.InfrastructureConfigurations, bus bus.Bus, metrics contracts.CatalogsMetrics) {
	routeBuilder.RegisterGroup("/api/v1", func(v1 *echo.Group) {
		group := v1.Group("/products")
		productEndpointBase := delivery.NewProductEndpointBase(infra, group, bus, metrics)

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
		searchProducts := searchingProductV1.NewSearchProductsEndpoint(productEndpointBase)
		searchProducts.MapRoute()

		// GetProductById
		getProductByIdEndpoint := gettingProductByIdV1.NewGetProductByIdEndpoint(productEndpointBase)
		getProductByIdEndpoint.MapRoute()

		// DeleteProduct
		deleteProductEndpoint := deletingProductV1.NewDeleteProductEndpoint(productEndpointBase)
		deleteProductEndpoint.MapRoute()
	},
	)
}
