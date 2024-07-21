package otelmetrics

import (
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

type config struct {
	metricsProvider metric.MeterProvider

	skipper middleware.Skipper

	namespace string

	serviceName string

	instrumentationName string

	// enableTotalMetric whether to enable a metric to count the total number of http requests.
	enableTotalMetric bool

	// enableDurMetric whether to enable a metric to track the duration of each request.
	enableDurMetric bool

	// enableDurMetric whether to enable a metric that tells the number of current in-flight requests.
	enableInFlightMetric bool
}

var defualtConfig = config{
	metricsProvider:      otel.GetMeterProvider(),
	enableTotalMetric:    true,
	enableDurMetric:      true,
	enableInFlightMetric: true,
	skipper:              middleware.DefaultSkipper,
	serviceName:          "application",
	instrumentationName:  "echo",
}

type Option interface {
	apply(*config)
}

type optionFunc func(*config)

func (o optionFunc) apply(c *config) {
	o(c)
}

// WithNamespace will set the metrics namespace that will be added to all metric configurations. It will be a prefix to each
// metric name. For example, if namespace is "myapp", then requests_total metric will be myapp_http_requests_total
// (after namespace there is also the subsystem prefix, "http" in this case).
func WithNamespace(v string) Option {
	return optionFunc(func(cfg *config) {
		if cfg.namespace != "" {
			cfg.namespace = v
		}
	})
}

func WithServiceName(v string) Option {
	return optionFunc(func(cfg *config) {
		if cfg.serviceName != "" {
			cfg.serviceName = v
		}
	})
}

func WithInstrumentationName(v string) Option {
	return optionFunc(func(cfg *config) {
		if cfg.instrumentationName != "" {
			cfg.instrumentationName = v
		}
	})
}

// WithSkipper specifies a skipper for allowing requests to skip generating spans.
func WithSkipper(skipper middleware.Skipper) Option {
	return optionFunc(func(cfg *config) {
		cfg.skipper = skipper
	})
}

// WithMeterProvider specifies a meter provider to use for creating a metrics.
// If none is specified, the global provider is used.
func WithMeterProvider(provider metric.MeterProvider) Option {
	return optionFunc(func(cfg *config) {
		if provider != nil {
			cfg.metricsProvider = provider
		}
	})
}

func WithInFlightMetric(enabled bool) Option {
	return optionFunc(func(cfg *config) {
		cfg.enableInFlightMetric = enabled
	})
}

func WithTotalMetric(enabled bool) Option {
	return optionFunc(func(cfg *config) {
		cfg.enableTotalMetric = enabled
	})
}

func WithDurMetric(enabled bool) Option {
	return optionFunc(func(cfg *config) {
		cfg.enableDurMetric = enabled
	})
}
