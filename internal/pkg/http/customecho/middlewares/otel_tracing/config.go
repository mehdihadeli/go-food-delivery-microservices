package oteltracing

import (
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// Ref: https://github.com/open-telemetry/opentelemetry-go-contrib/blob/main/instrumentation/github.com/labstack/echo/otelecho/echo.go

// config is used to configure the mux middleware.
type config struct {
	tracerProvider      oteltrace.TracerProvider
	propagators         propagation.TextMapPropagator
	skipper             middleware.Skipper
	instrumentationName string
	serviceName         string
}

// Option specifies instrumentation configuration options.
type Option interface {
	apply(*config)
}

type optionFunc func(*config)

func (o optionFunc) apply(c *config) {
	o(c)
}

var defualtConfig = config{
	tracerProvider:      otel.GetTracerProvider(),
	propagators:         otel.GetTextMapPropagator(),
	skipper:             middleware.DefaultSkipper,
	instrumentationName: "echo",
	serviceName:         "app",
}

// WithPropagators specifies propagators to use for extracting
// information from the HTTP requests. If none are specified, global
// ones will be used.
func WithPropagators(propagators propagation.TextMapPropagator) Option {
	return optionFunc(func(cfg *config) {
		if propagators != nil {
			cfg.propagators = propagators
		}
	})
}

// WithTracerProvider specifies a tracer provider to use for creating a tracer.
// If none is specified, the global provider is used.
func WithTracerProvider(provider oteltrace.TracerProvider) Option {
	return optionFunc(func(cfg *config) {
		if provider != nil {
			cfg.tracerProvider = provider
		}
	})
}

// WithSkipper specifies a skipper for allowing requests to skip generating spans.
func WithSkipper(skipper middleware.Skipper) Option {
	return optionFunc(func(cfg *config) {
		cfg.skipper = skipper
	})
}

func WithInstrumentationName(v string) Option {
	return optionFunc(func(cfg *config) {
		if cfg.instrumentationName != "" {
			cfg.instrumentationName = v
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
