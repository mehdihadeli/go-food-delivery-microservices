package delivery

import (
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/configurations/infrastructure"
)

type ProductEndpointBase struct {
	*infrastructure.InfrastructureConfigurations
	ProductMediator *mediatr.Mediator
	ProductsGroup   *echo.Group
}
