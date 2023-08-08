package products

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/data/repositories"
	createProductV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/creating_product/v1/endpoints"
	deleteProductV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/deleting_product/v1/endpoints"
	getProductByIdV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/getting_product_by_id/v1/endpoints"
	getProductsV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/getting_products/v1/endpoints"
	searchProductsV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/searching_product/v1/endpoints"
	updateProductsV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/updating_product/v1/endpoints"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/shared/grpc"

	customEcho "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/custom_echo"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/web/route"
)

var Module = fx.Module(
	"productsfx",

	// Other provides
	fx.Provide(repositories.NewPostgresProductRepository),
	fx.Provide(grpc.NewProductGrpcService),

	fx.Provide(fx.Annotate(func(catalogsServer customEcho.EchoHttpServer) *echo.Group {
		var g *echo.Group
		catalogsServer.RouteBuilder().RegisterGroupFunc("/api/v1", func(v1 *echo.Group) {
			group := v1.Group("/products")
			g = group
		})

		return g
	}, fx.ResultTags(`name:"product-echo-group"`))),

	fx.Provide(
		route.AsRoute(createProductV1.NewCreteProductEndpoint, "product-routes"),
		route.AsRoute(updateProductsV1.NewUpdateProductEndpoint, "product-routes"),
		route.AsRoute(getProductsV1.NewGetProductsEndpoint, "product-routes"),
		route.AsRoute(searchProductsV1.NewSearchProductsEndpoint, "product-routes"),
		route.AsRoute(getProductByIdV1.NewGetProductByIdEndpoint, "product-routes"),
		route.AsRoute(deleteProductV1.NewDeleteProductEndpoint, "product-routes"),
	),
)
