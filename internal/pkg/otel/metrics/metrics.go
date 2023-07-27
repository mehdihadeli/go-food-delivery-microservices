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
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/config"
)

type OtelMetrics struct {
	Config *config.OpenTelemetryOptions
	Logger logger.Logger
	Meter  api.Meter
	Echo   *echo.Echo
}

// NewOtelMetrics adds otel metrics
// https://github.com/open-telemetry/opentelemetry-go/blob/main/example/prometheus/main.go
func NewOtelMetrics(
	config *config.OpenTelemetryOptions,
	logger logger.Logger,
) (*OtelMetrics, error) {
	if config == nil {
		return nil, errors.New("metrics config can't be nil")
	}

	e := echo.New()
	e.HideBanner = false

	// The exporter embeds a default OpenTelemetry Reader and
	// implements prometheus.Collector, allowing it to be used as
	// both a Reader and Collector.
	exporter, err := prometheus.New()
	if err != nil {
		logger.Fatal(err)
	}
	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	meter := provider.Meter(
		"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/metrics",
	)

	return &OtelMetrics{Config: config, Meter: meter, Logger: logger, Echo: e}, nil
}

func (o *OtelMetrics) Run() error {
	o.Echo.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         1 << 10, // 1 KB
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))

	var metricsPath string
	if o.Config.OTelMetricsOptions.MetricsRoutePath == "" {
		metricsPath = "/metrics"
	} else {
		metricsPath = o.Config.OTelMetricsOptions.MetricsRoutePath
	}

	o.Echo.GET(metricsPath, echo.WrapHandler(promhttp.Handler()))
	o.Logger.Infof("serving metrics at localhost:%s/metrics", o.Config.OTelMetricsOptions.Port)
	err := o.Echo.Start(o.Config.OTelMetricsOptions.Port)

	return err
}

func (o *OtelMetrics) GracefulShutdown(ctx context.Context) error {
	err := o.Echo.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}
