package infrastructure

import customMiddlewares "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/web/custom_middlewares"

func (ic *infrastructureConfigurator) configMiddlewares(metrics *OrdersServiceMetrics) {

	ic.echoServer.SetupDefaultMiddlewares()

	middlewares := customMiddlewares.NewCustomMiddlewares(ic.log, ic.cfg, getHttpMetricsCb(metrics))
	ic.echoServer.AddMiddlewares(middlewares.RequestLoggerMiddleware, middlewares.RequestMetricsMiddleware)
}

func getHttpMetricsCb(metrics *OrdersServiceMetrics) func(err error) {
	return func(err error) {
		if err != nil {
			metrics.ErrorHttpRequests.Inc()
		} else {
			metrics.SuccessHttpRequests.Inc()
		}
	}
}

func getGrpcMetricsCb(metrics *OrdersServiceMetrics) func(err error) {
	return func(err error) {
		if err != nil {
			metrics.ErrorGrpcRequests.Inc()
		} else {
			metrics.SuccessGrpcRequests.Inc()
		}
	}
}
