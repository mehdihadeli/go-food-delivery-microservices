package metrics

import (
	"context"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/customecho/contracts"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger"

	"go.opentelemetry.io/contrib/instrumentation/host"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/fx"
)

var (
	// Module provided to fxlog
	// https://uber-go.github.io/fx/modules.html
	Module = fx.Module( //nolint:gochecknoglobals
		"otelmetrixfx",
		metricsProviders,
		metricsInvokes,
	)

	metricsProviders = fx.Options(fx.Provide( //nolint:gochecknoglobals
		ProvideMetricsConfig,
		NewOtelMetrics,
		fx.Annotate(
			provideMeter,
			fx.ParamTags(`optional:"true"`),
			fx.As(new(AppMetrics)),
			fx.As(new(metric.Meter))),
	))

	metricsInvokes = fx.Options( //nolint:gochecknoglobals
		fx.Invoke(registerHooks),
		fx.Invoke(func(m *OtelMetrics, server contracts.EchoHttpServer) {
			m.RegisterMetricsEndpoint(server)
		}),
	)
)

func provideMeter(otelMetrics *OtelMetrics) AppMetrics {
	return otelMetrics.appMetrics
}

// we don't want to register any dependencies here, its func body should execute always even we don't request for that, so we should use `invoke`
func registerHooks(
	lc fx.Lifecycle,
	metrics *OtelMetrics,
	logger logger.Logger,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if metrics.appMetrics == nil {
				return nil
			}

			if metrics.config.EnableHostMetrics {
				// https://github.com/open-telemetry/opentelemetry-go-contrib/tree/main/instrumentation/host
				// we changed default meter provider in metrics setup
				logger.Info("Starting host instrumentation:")
				err := host.Start(
					host.WithMeterProvider(otel.GetMeterProvider()),
				)
				if err != nil {
					logger.Errorf(
						"error starting host instrumentation: %s",
						err,
					)
				}
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {
			if metrics.appMetrics == nil {
				return nil
			}

			if err := metrics.Shutdown(ctx); err != nil {
				logger.Errorf(
					"error in shutting down metrics provider: %v",
					err,
				)
			} else {
				logger.Info("metrics provider shutdown gracefully")
			}

			return nil
		},
	})
}
