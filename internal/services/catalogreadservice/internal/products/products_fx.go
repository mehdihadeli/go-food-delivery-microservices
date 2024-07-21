package products

import (
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/web/route"
	customEcho "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/customecho"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogreadservice/internal/products/data/repositories"
	getProductByIdV1 "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogreadservice/internal/products/features/get_product_by_id/v1/endpoints"
	getProductsV1 "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogreadservice/internal/products/features/getting_products/v1/endpoints"
	searchProductV1 "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogreadservice/internal/products/features/searching_products/v1/endpoints"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"productsfx",

	// Other provides
	fx.Provide(repositories.NewRedisProductRepository),
	fx.Provide(repositories.NewMongoProductRepository),

	fx.Provide(fx.Annotate(func(catalogsServer customEcho.EchoHttpServer) *echo.Group {
		var g *echo.Group
		catalogsServer.RouteBuilder().RegisterGroupFunc("/api/v1", func(v1 *echo.Group) {
			group := v1.Group("/products")
			g = group
		})

		return g
	}, fx.ResultTags(`name:"product-echo-group"`))),

	fx.Provide(
		route.AsRoute(getProductsV1.NewGetProductsEndpoint, "product-routes"),
		route.AsRoute(searchProductV1.NewSearchProductsEndpoint, "product-routes"),
		route.AsRoute(getProductByIdV1.NewGetProductByIdEndpoint, "product-routes"),
	),
)
