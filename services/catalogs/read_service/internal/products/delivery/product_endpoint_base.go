package delivery

import (
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/contracts"
)

type ProductEndpointBase struct {
	contracts.InfrastructureConfigurations
	ProductsGroup *echo.Group
}

func NewProductEndpointBase(infra contracts.InfrastructureConfigurations, productsGroup *echo.Group) *ProductEndpointBase {
	return &ProductEndpointBase{InfrastructureConfigurations: infra, ProductsGroup: productsGroup}
}
