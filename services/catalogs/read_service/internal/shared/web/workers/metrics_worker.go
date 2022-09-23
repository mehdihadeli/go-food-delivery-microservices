package workers

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/constants"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/web"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/configurations/infrastructure"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewMetricsWorker(infra *infrastructure.InfrastructureConfigurations) web.Worker {
	metricsServer := echo.New()
	return web.NewBackgroundWorker(func(ctx context.Context) error {
		metricsServer.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
			StackSize:         constants.StackSize,
			DisablePrintStack: true,
			DisableStackAll:   true,
		}))
		metricsServer.GET(infra.Cfg.Probes.PrometheusPath, echo.WrapHandler(promhttp.Handler()))
		infra.Log.Infof("Metrics server is running on port: %s", infra.Cfg.Probes.PrometheusPort)
		if err := metricsServer.Start(infra.Cfg.Probes.PrometheusPort); err != nil {
			infra.Log.Errorf("metricsServer.Start: %v", err)
			return err
		}
		return nil
	}, func(ctx context.Context) error {
		return metricsServer.Shutdown(ctx)
	})
}
