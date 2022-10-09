package delivery

import (
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/bus"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/contracts"
)

type OrderEndpointBase struct {
	contracts.InfrastructureConfigurations
	OrdersGroup   *echo.Group
	OrdersMetrics contracts.OrdersMetrics
	Bus           bus.Bus
}

func NewOrderEndpointBase(infra contracts.InfrastructureConfigurations, ordersGroup *echo.Group, bus bus.Bus, ordersMetrics contracts.OrdersMetrics) *OrderEndpointBase {
	return &OrderEndpointBase{OrdersGroup: ordersGroup, InfrastructureConfigurations: infra, Bus: bus, OrdersMetrics: ordersMetrics}
}
