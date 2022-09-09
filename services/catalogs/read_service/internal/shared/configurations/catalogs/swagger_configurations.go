package catalogs

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func (c *catalogsServiceConfigurator) configSwagger() {
	//https://github.com/swaggo/swag#how-to-use-it-with-gin
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Title = "Catalogs Read-Service Api"
	docs.SwaggerInfo.Description = "Catalogs Read-Service Api."

	c.echoServer.GetEchoInstance().GET("/swagger/*", echoSwagger.WrapHandler)
}
