package catalogs

import (
	"github.com/labstack/echo/v4"
	customEcho "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/custom_echo"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/docs"
)

func (ic *catalogsServiceConfigurator) configSwagger(routeBuilder *customEcho.RouteBuilder) {
	//https://github.com/swaggo/swag#how-to-use-it-with-gin
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Title = "Catalogs Read-Service Api"
	docs.SwaggerInfo.Description = "Catalogs Read-Service Api."

	routeBuilder.RegisterRoutes(func(e *echo.Echo) {
		e.GET("/swagger/*", echoSwagger.WrapHandler)
	})
}
