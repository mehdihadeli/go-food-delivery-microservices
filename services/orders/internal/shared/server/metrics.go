package server

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/constants"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (s *Server) RunMetrics(cancel context.CancelFunc) {
	metricsServer := echo.New()
	go func() {
		metricsServer.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
			StackSize:         constants.StackSize,
			DisablePrintStack: true,
			DisableStackAll:   true,
		}))
		metricsServer.GET(s.cfg.Probes.PrometheusPath, echo.WrapHandler(promhttp.Handler()))
		s.log.Infof("Metrics server is running on port: %s", s.cfg.Probes.PrometheusPort)
		if err := metricsServer.Start(s.cfg.Probes.PrometheusPort); err != nil {
			s.log.Errorf("metricsServer.Start: %v", err)
			cancel()
		}
	}()
}
