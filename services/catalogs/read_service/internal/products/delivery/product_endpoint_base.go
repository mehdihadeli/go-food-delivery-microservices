package delivery

import (
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/contracts"
)

type ProductEndpointBase struct {
	contracts.InfrastructureConfiguration
	ProductsGroup *echo.Group
}

func NewProductEndpointBase(infra contracts.InfrastructureConfiguration, productsGroup *echo.Group) *ProductEndpointBase {
	return &ProductEndpointBase{InfrastructureConfiguration: infra, ProductsGroup: productsGroup}
}
