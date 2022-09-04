package delivery

import (
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/configurations/infrastructure"
)

type ProductEndpointBase struct {
	*infrastructure.InfrastructureConfiguration
	ProductsGroup *echo.Group
}

func NewProductEndpointBase(infra *infrastructure.InfrastructureConfiguration, productsGroup *echo.Group) *ProductEndpointBase {
	return &ProductEndpointBase{ProductsGroup: productsGroup, InfrastructureConfiguration: infra}
}
