package infrastructure

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func (ic *infrastructureConfigurator) configSwagger() {
	//https://github.com/swaggo/swag#how-to-use-it-with-gin
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Title = "Catalogs Write-Service Api"
	docs.SwaggerInfo.Description = "Catalogs Write-Service Api."

	ic.echoServer.GetEchoInstance().GET("/swagger/*", echoSwagger.WrapHandler)
}
