package configurations

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func (ic *infrastructureConfigurator) configSwagger() {
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Title = "Catalogs Service Api"
	docs.SwaggerInfo.Description = "Catalogs Service Api."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api/v1"

	ic.echo.GET("/swagger/*", echoSwagger.WrapHandler)
}

