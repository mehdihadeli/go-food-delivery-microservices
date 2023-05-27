package metrics

import (
    "context"

    "emperror.dev/errors"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "go.opentelemetry.io/otel/exporters/prometheus"
    api "go.opentelemetry.io/otel/metric"
    "go.opentelemetry.io/otel/sdk/metric"

    "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
)

// AddOtelMetrics adds otel metrics
// https://github.com/open-telemetry/opentelemetry-go/blob/main/example/prometheus/main.go
func AddOtelMetrics(ctx context.Context, config *OTelMetricsConfig, logger logger.Logger) (api.Meter, error) {
	if config == nil {
		return nil, errors.New("metrics config can't be nil")
	}

	// The exporter embeds a default OpenTelemetry Reader and
	// implements prometheus.Collector, allowing it to be used as
	// both a Reader and Collector.
	exporter, err := prometheus.New()
	if err != nil {
		logger.Fatal(err)
	}
	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	meter := provider.Meter("github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/metrics")

	// Start the prometheus HTTP server and pass the exporter Collector to it
	go serveMetrics(ctx, logger, config)

	return meter, nil
}

func serveMetrics(ctx context.Context, logger logger.Logger, config *OTelMetricsConfig) {
	echoServer := echo.New()
	echoServer.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         1 << 10, // 1 KB
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))

	go func() {
		for {
			select {
			case <-ctx.Done():
				logger.Infof("%s is shutting down, Metrics Http PORT: {%s}", config.Name, config.Port)
				if err := echoServer.Shutdown(ctx); err != nil {
					logger.Errorf("(Shutdown) err: {%v}", err)
				}
				return
			}
		}
	}()

	var metricsPath string
	if config.MetricsRoutePath == "" {
		metricsPath = "/metrics"
	} else {
		metricsPath = config.MetricsRoutePath
	}

	echoServer.GET(metricsPath, echo.WrapHandler(promhttp.Handler()))
	logger.Infof("serving metrics at localhost:%s/metrics", config.Port)
	err := echoServer.Start(config.Port)
	if err != nil {
		logger.Errorf("error serving http: %v", err)
		return
	}
}
