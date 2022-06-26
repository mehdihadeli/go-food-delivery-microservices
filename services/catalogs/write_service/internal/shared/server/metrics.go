package server

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	catalog_constants "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/constants"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (s *Server) RunMetrics(cancel context.CancelFunc) {
	metricsServer := echo.New()
	go func() {
		metricsServer.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
			StackSize:         catalog_constants.StackSize,
			DisablePrintStack: true,
			DisableStackAll:   true,
		}))
		metricsServer.GET(s.Cfg.Probes.PrometheusPath, echo.WrapHandler(promhttp.Handler()))
		s.Log.Infof("Metrics server is running on port: %s", s.Cfg.Probes.PrometheusPort)
		if err := metricsServer.Start(s.Cfg.Probes.PrometheusPort); err != nil {
			s.Log.Errorf("metricsServer.Start: %v", err)
			cancel()
		}
	}()
}
