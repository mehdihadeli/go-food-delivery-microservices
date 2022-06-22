package server

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	catalog_constants "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/shared/constants"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (i *Server) RunMetrics(cancel context.CancelFunc) {
	metricsServer := echo.New()
	go func() {
		metricsServer.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
			StackSize:         catalog_constants.StackSize,
			DisablePrintStack: true,
			DisableStackAll:   true,
		}))
		metricsServer.GET(i.Cfg.Probes.PrometheusPath, echo.WrapHandler(promhttp.Handler()))
		i.Log.Infof("Metrics server is running on port: %s", i.Cfg.Probes.PrometheusPort)
		if err := metricsServer.Start(i.Cfg.Probes.PrometheusPort); err != nil {
			i.Log.Errorf("metricsServer.Start: %v", err)
			cancel()
		}
	}()
}
