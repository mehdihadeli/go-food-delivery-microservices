package otel

import (
	"context"
	"errors"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/metric"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/metrics"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing"
)

// Module provided to fxlog
// https://uber-go.github.io/fx/modules.html
var Module = fx.Module(
	"otelfx",
	fx.Provide(
		config.ProvideOtelConfig,
		metrics.NewOtelMetrics,
		fx.Annotate(
			provideMeter,
			fx.As(new(metric.Meter))),
		tracing.NewOtelTracing,
	),
	fx.Invoke(registerHooks),
)

func provideMeter(otelMetrics *metrics.OtelMetrics) metric.Meter {
	return otelMetrics.Meter
}

// we don't want to register any dependencies here, its func body should execute always even we don't request for that, so we should use `invoke`
func registerHooks(
	lc fx.Lifecycle,
	metrics *metrics.OtelMetrics,
	tracerProvider *tracesdk.TracerProvider,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Start server in a separate goroutine, this way when the server is shutdown "s.e.Start" will
			// return promptly, and the call to "s.e.Shutdown" is the one that will wait for all other
			// resources to be properly freed. If it was the other way around, the application would just
			// exit without gracefully shutting down the server.
			// For more details: https://medium.com/@momchil.dev/proper-http-shutdown-in-go-bd3bfaade0f2
			go func() {
				if err := metrics.Run(ctx); !errors.Is(
					err,
					http.ErrServerClosed,
				) {
					metrics.Logger.Fatalf("(s.RunHttpServer) error in running server: {%v}", err)
				}
			}()
			metrics.Logger.Infof(
				"Metrics server %s is listening on Host:{%s} Http PORT: {%s}",
				metrics.Config.OTelMetricsOptions.Name,
				metrics.Config.OTelMetricsOptions.Host,
				metrics.Config.OTelMetricsOptions.Port,
			)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
			defer cancel()

			if err := tracerProvider.Shutdown(ctx); err != nil {
				metrics.Logger.Errorf("error in shutting down trace provider: %v", err)
			} else {
				metrics.Logger.Info("trace provider shutdown gracefully")
			}

			if err := metrics.GracefulShutdown(ctx); err != nil {
				metrics.Logger.Errorf("error shutting down metrics server: %v", err)
			} else {
				metrics.Logger.Info("metrics server shutdown gracefully")
			}
			return nil
		},
	})
}
