package configurations

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared"
	catalog_constants "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/constants"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/web/middlewares"
	"strings"
)

func (ic *infrastructureConfigurator) configMiddlewares() {
	ic.echo.HideBanner = false

	ic.echo.HTTPErrorHandler = middlewares.ProblemHandler

	middlewareManager := middlewares.NewMiddlewareManager(ic.log, ic.cfg, getHttpMetricsCb(infrastructure.Metrics))

	ic.echo.Use(middlewareManager.RequestLoggerMiddleware)
	ic.echo.Use(middlewareManager.RequestMetricsMiddleware)

	ic.echo.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         catalog_constants.StackSize,
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))
	ic.echo.Use(middleware.RequestID())
	ic.echo.Use(middleware.Logger())
	ic.echo.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: catalog_constants.GzipLevel,
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "swagger")
		},
	}))

	ic.echo.Use(middleware.BodyLimit(catalog_constants.BodyLimit))
}

func getHttpMetricsCb(metrics *shared.CatalogsServiceMetrics) func(err error) {
	return func(err error) {
		if err != nil {
			metrics.ErrorHttpRequests.Inc()
		} else {
			metrics.SuccessHttpRequests.Inc()
		}
	}
}

func getGrpcMetricsCb(metrics *shared.CatalogsServiceMetrics) func(err error) {
	return func(err error) {
		if err != nil {
			metrics.ErrorGrpcRequests.Inc()
		} else {
			metrics.SuccessGrpcRequests.Inc()
		}
	}
}
