package configurations

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/delivery"
	creatingOrderV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/creating_order/endpoints/v1"
	gettingOrderByIdV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/getting_order_by_id/endpoints/v1"

	"github.com/labstack/echo/v4"
)

func (c *ordersModuleConfigurator) configEndpoints(ctx context.Context, group *echo.Group) {

	orderEndpointBase := &delivery.OrderEndpointBase{
		OrdersGroup:                 group,
		InfrastructureConfiguration: c.InfrastructureConfiguration,
	}

	// CreateNewOrder
	createProductEndpoint := creatingOrderV1.NewCreteOrderEndpoint(orderEndpointBase)
	createProductEndpoint.MapRoute()

	// GetOrderByID
	getOrderByIdEndpoint := gettingOrderByIdV1.NewGetOrderByIdEndpoint(orderEndpointBase)
	getOrderByIdEndpoint.MapRoute()
}
