package orders

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/configurations/infrastructure"
	customMiddlewares "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/web/custom_middlewares"
)

func (c *ordersServiceConfigurator) configMiddlewares(metrics *infrastructure.OrdersServiceMetrics) {
	c.echoServer.SetupDefaultMiddlewares()

	middlewares := customMiddlewares.NewCustomMiddlewares(c.Log, c.Cfg, getHttpMetricsCb(metrics))
	c.echoServer.AddMiddlewares(middlewares.RequestMetricsMiddleware)
}

func getHttpMetricsCb(metrics *infrastructure.OrdersServiceMetrics) func(err error) {
	return func(err error) {
		if err != nil {
			metrics.ErrorHttpRequests.Inc()
		} else {
			metrics.SuccessHttpRequests.Inc()
		}
	}
}

func getGrpcMetricsCb(metrics *infrastructure.OrdersServiceMetrics) func(err error) {
	return func(err error) {
		if err != nil {
			metrics.ErrorGrpcRequests.Inc()
		} else {
			metrics.SuccessGrpcRequests.Inc()
		}
	}
}
