package configurations

import (
	"context"
	"github.com/labstack/echo/v4"
	customEcho "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/custom_echo"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/delivery"
	creatingOrderV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/creating_order/endpoints/v1"
	gettingOrderByIdV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/getting_order_by_id/endpoints/v1"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/configurations/infrastructure"
)

func (c *ordersModuleConfigurator) configEndpoints(ctx context.Context) {
	configV1Endpoints(c.EchoServer, c.InfrastructureConfiguration, ctx)
}

func configV1Endpoints(echoServer customEcho.EchoHttpServer, infra *infrastructure.InfrastructureConfiguration, ctx context.Context) {
	echoServer.ConfigGroup("/api/v1", func(v1 *echo.Group) {
		ordersGroup := v1.Group("/orders")

		orderEndpointBase := &delivery.OrderEndpointBase{
			OrdersGroup:                 ordersGroup,
			InfrastructureConfiguration: infra,
		}

		// CreateNewOrder
		createProductEndpoint := creatingOrderV1.NewCreteOrderEndpoint(orderEndpointBase)
		createProductEndpoint.MapRoute()

		// GetOrderByID
		getOrderByIdEndpoint := gettingOrderByIdV1.NewGetOrderByIdEndpoint(orderEndpointBase)
		getOrderByIdEndpoint.MapRoute()
	})
}
