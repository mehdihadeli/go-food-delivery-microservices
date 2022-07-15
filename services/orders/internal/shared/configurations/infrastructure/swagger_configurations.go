package infrastructure

import (
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/swaggo/swag/example/basic/docs"
)

func (ic *infrastructureConfigurator) configSwagger() {
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Title = "Orders Service Api"
	docs.SwaggerInfo.Description = "Orders Service Api."

	ic.echo.GET("/swagger/*", echoSwagger.WrapHandler)
}
