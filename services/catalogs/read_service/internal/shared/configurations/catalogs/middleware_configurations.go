package catalogs

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/configurations/infrastructure"
	customMiddlewares "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/web/middlewares"
)

func (c *catalogsServiceConfigurator) configMiddlewares(metrics *infrastructure.CatalogsServiceMetrics) {
	c.echoServer.SetupDefaultMiddlewares()

	middlewares := customMiddlewares.NewCustomMiddlewares(c.Log, c.Cfg, getHttpMetricsCb(metrics))
	c.echoServer.AddMiddlewares(middlewares.RequestMetricsMiddleware)
}

func getHttpMetricsCb(metrics *infrastructure.CatalogsServiceMetrics) func(err error) {
	return func(err error) {
		if err != nil {
			metrics.ErrorHttpRequests.Inc()
		} else {
			metrics.SuccessHttpRequests.Inc()
		}
	}
}

func getGrpcMetricsCb(metrics *infrastructure.CatalogsServiceMetrics) func(err error) {
	return func(err error) {
		if err != nil {
			metrics.ErrorGrpcRequests.Inc()
		} else {
			metrics.SuccessGrpcRequests.Inc()
		}
	}
}
