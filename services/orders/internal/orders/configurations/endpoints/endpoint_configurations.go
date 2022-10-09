package endpoints

import (
	"context"
	"github.com/labstack/echo/v4"
	customEcho "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/custom_echo"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/bus"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/delivery"
	creatingOrderV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/creating_order/endpoints/v1"
	gettingOrderByIdV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/getting_order_by_id/endpoints/v1"
	gettingOrdersV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/getting_orders/endpoints/v1"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/contracts"
)

func ConfigOrdersEndpoints(ctx context.Context, routeBuilder *customEcho.RouteBuilder, infra contracts.InfrastructureConfigurations, bus bus.Bus, metrics contracts.OrdersMetrics) {
	configV1Endpoints(ctx, routeBuilder, infra, bus, metrics)
}

func configV1Endpoints(ctx context.Context, routeBuilder *customEcho.RouteBuilder, infra contracts.InfrastructureConfigurations, bus bus.Bus, metrics contracts.OrdersMetrics) {
	routeBuilder.RegisterGroup("/api/v1", func(v1 *echo.Group) {
		ordersGroup := v1.Group("/orders")

		orderEndpointBase := delivery.NewOrderEndpointBase(infra, ordersGroup, bus, metrics)

		// CreateNewOrder
		createOrderEndpoint := creatingOrderV1.NewCreteOrderEndpoint(orderEndpointBase)
		createOrderEndpoint.MapRoute()

		// GetOrderByID
		getOrderByIdEndpoint := gettingOrderByIdV1.NewGetOrderByIdEndpoint(orderEndpointBase)
		getOrderByIdEndpoint.MapRoute()

		// GetOrders
		getOrders := gettingOrdersV1.NewGetOrdersEndpoint(orderEndpointBase)
		getOrders.MapRoute()
	})
}
