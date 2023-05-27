package endpoints

import (
	"context"

	"github.com/labstack/echo/v4"

	customEcho "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/custom_echo"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/bus"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/delivery"
	createProductV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/features/creating_product/v1/endpoints"
	deleteProductV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/features/deleting_product/v1/endpoints"
	getProductByIdV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/features/getting_product_by_id/v1/endpoints"
	getProductsV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/features/getting_products/v1/endpoints"
	searchProductsV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/features/searching_product/v1/endpoints"
	updateProductsV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/features/updating_product/v1/endpoints"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/shared/contracts"
)

func ConfigProductsEndpoints(ctx context.Context, routeBuilder *customEcho.RouteBuilder, infra *contracts.InfrastructureConfigurations, bus bus.Bus, metrics *contracts.CatalogsMetrics) {
	configV1Endpoints(ctx, routeBuilder, infra, bus, metrics)
}

func configV1Endpoints(ctx context.Context, routeBuilder *customEcho.RouteBuilder, infra *contracts.InfrastructureConfigurations, bus bus.Bus, metrics *contracts.CatalogsMetrics) {
	routeBuilder.RegisterGroupFunc("/api/v1", func(v1 *echo.Group) {
		group := v1.Group("/products")
		productEndpointBase := delivery.NewProductEndpointBase(infra, group, bus, metrics)

		// CreateNewProduct
		createProductEndpoint := createProductV1.NewCreteProductEndpoint(productEndpointBase)
		createProductEndpoint.MapRoute()

		// UpdateProduct
		updateProductEndpoint := updateProductsV1.NewUpdateProductEndpoint(productEndpointBase)
		updateProductEndpoint.MapRoute()

		// GetProducts
		getProductsEndpoint := getProductsV1.NewGetProductsEndpoint(productEndpointBase)
		getProductsEndpoint.MapRoute()

		// SearchProducts
		searchProducts := searchProductsV1.NewSearchProductsEndpoint(productEndpointBase)
		searchProducts.MapRoute()

		// GetProductById
		getProductByIdEndpoint := getProductByIdV1.NewGetProductByIdEndpoint(productEndpointBase)
		getProductByIdEndpoint.MapRoute()

		// DeleteProduct
		deleteProductEndpoint := deleteProductV1.NewDeleteProductEndpoint(productEndpointBase)
		deleteProductEndpoint.MapRoute()
	},
	)
}
