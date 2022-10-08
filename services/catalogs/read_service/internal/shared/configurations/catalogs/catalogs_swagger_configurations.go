package catalogs

import (
	"github.com/labstack/echo/v4"
	customEcho "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/custom_echo"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
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
