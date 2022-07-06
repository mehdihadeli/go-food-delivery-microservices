package delivery

import (
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/configurations/infrastructure"
)

type ProductEndpointBase struct {
	*infrastructure.InfrastructureConfiguration
	ProductMediator *mediatr.Mediator
	ProductsGroup   *echo.Group
}
