package infrastructure

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func (ic *infrastructureConfigurator) configSwagger() {
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Title = "Catalogs Read-Service Api"
	docs.SwaggerInfo.Description = "Catalogs Read-Service Api."

	ic.echo.GET("/swagger/*", echoSwagger.WrapHandler)
}
