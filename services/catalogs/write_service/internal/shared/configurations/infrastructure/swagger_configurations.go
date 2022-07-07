package infrastructure

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func (ic *infrastructureConfigurator) configSwagger() {
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Title = "Catalogs Write-Service Api"
	docs.SwaggerInfo.Description = "Catalogs Write-Service Api."

	ic.echo.GET("/swagger/*", echoSwagger.WrapHandler)
}
