package delivery

import (
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/configurations/infrastructure"
)

type OrderEndpointBase struct {
	*infrastructure.InfrastructureConfiguration
	OrdersGroup *echo.Group
}

func NewOrderEndpointBase(infra *infrastructure.InfrastructureConfiguration, ordersGroup *echo.Group) *OrderEndpointBase {
	return &OrderEndpointBase{OrdersGroup: ordersGroup, InfrastructureConfiguration: infra}
}
