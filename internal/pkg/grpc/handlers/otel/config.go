package otel

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// ref: https://github.com/bakins/otel-grpc-statshandler/blob/main/statshandler.go

// Option applies an option value when creating a Handler.
type Option interface {
	apply(*config)
}

type optionFunc func(*config)

func (f optionFunc) apply(c *config) {
	f(c)
}

type config struct {
	metricsProvider     metric.MeterProvider
	tracerProvider      trace.TracerProvider
	propagator          propagation.TextMapPropagator
	Namespace           string
	serviceName         string
	instrumentationName string
}

var defualtConfig = config{
	metricsProvider:     otel.GetMeterProvider(),
	tracerProvider:      otel.GetTracerProvider(),
	propagator:          otel.GetTextMapPropagator(),
	serviceName:         "application",
	instrumentationName: "grpc-otel",
}

func WithMeterProvider(m metric.MeterProvider) Option {
	return optionFunc(func(c *config) {
		c.metricsProvider = m
	})
}

func WithTraceProvider(t trace.TracerProvider) Option {
	return optionFunc(func(c *config) {
		c.tracerProvider = t
	})
}

func WithPropagators(p propagation.TextMapPropagator) Option {
	return optionFunc(func(c *config) {
		c.propagator = p
	})
}

func SetNamespace(v string) Option {
	return optionFunc(func(cfg *config) {
		if cfg.Namespace != "" {
			cfg.Namespace = v
		}
	})
}

func SetServiceName(v string) Option {
	return optionFunc(func(cfg *config) {
		if cfg.serviceName != "" {
			cfg.serviceName = v
		}
	})
}

func SetInstrumentationName(v string) Option {
	return optionFunc(func(cfg *config) {
		if cfg.instrumentationName != "" {
			cfg.instrumentationName = v
		}
	})
}
