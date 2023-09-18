package otel

import (
	"context"
	"net/http"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/metrics"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing"

	"emperror.dev/errors"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

var (
	// Module provided to fxlog
	// https://uber-go.github.io/fx/modules.html
	Module = fx.Module( //nolint:gochecknoglobals
		"otelfx",
		otelProviders,
		otelInvokes,
	)

	otelProviders = fx.Options(fx.Provide( //nolint:gochecknoglobals
		config.ProvideOtelConfig,
		metrics.NewOtelMetrics,
		tracing.NewOtelTracing,
		fx.Annotate(
			provideMeter,
			fx.As(new(metric.Meter))),
		fx.Annotate(
			provideTracer,
			fx.As(new(tracing.AppTracer)),
			fx.As(new(trace.Tracer)),
		),
	))

	otelInvokes = fx.Options(fx.Invoke(registerHooks)) //nolint:gochecknoglobals
)

func provideMeter(otelMetrics *metrics.OtelMetrics) metric.Meter {
	return otelMetrics.Meter
}

func provideTracer(tracingOtel *tracing.TracingOpenTelemetry) tracing.AppTracer {
	return tracingOtel.AppTracer
}

// we don't want to register any dependencies here, its func body should execute always even we don't request for that, so we should use `invoke`
func registerHooks(
	lc fx.Lifecycle,
	metrics *metrics.OtelMetrics,
	logger logger.Logger,
	tracingOtel *tracing.TracingOpenTelemetry,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if metrics.Meter == nil {
				return nil
			}

			go func() {
				// https://medium.com/@mokiat/proper-http-shutdown-in-go-bd3bfaade0f2
				// When Shutdown is called, Serve, ListenAndServe, and ListenAndServeTLS immediately return ErrServerClosed. Make sure the program doesn’t exit and waits instead for Shutdown to return.
				if err := metrics.Run(); !errors.Is(err, http.ErrServerClosed) {
					// do a fatal for running OnStop hook
					logger.Fatalf(
						"(metrics.RunHttpServer) error in running metrics server: {%v}",
						err,
					)
				}
			}()
			logger.Infof(
				"Metrics server %s is listening on Host:{%s} Http PORT: {%s}",
				metrics.Config.OTelMetricsOptions.Name,
				metrics.Config.OTelMetricsOptions.Host,
				metrics.Config.OTelMetricsOptions.Port,
			)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			if err := tracingOtel.TracerProvider.Shutdown(ctx); err != nil {
				logger.Errorf("error in shutting down trace provider: %v", err)
			} else {
				logger.Info("trace provider shutdown gracefully")
			}

			if metrics.Meter == nil {
				return nil
			}
			// https://github.com/uber-go/fx/blob/v1.20.0/app.go#L573
			// this ctx is just for stopping callbacks or OnStop callbacks, and it has short timeout 15s, and it is not alive in whole lifetime app
			// https://medium.com/@mokiat/proper-http-shutdown-in-go-bd3bfaade0f2
			// When Shutdown is called, Serve, ListenAndServe, and ListenAndServeTLS immediately return ErrServerClosed. Make sure the program doesn’t exit and waits instead for Shutdown to return.
			if err := metrics.GracefulShutdown(ctx); err != nil {
				logger.Errorf("error shutting down metrics server: %v", err)
			} else {
				logger.Info("metrics server shutdown gracefully")
			}
			return nil
		},
	})
}
