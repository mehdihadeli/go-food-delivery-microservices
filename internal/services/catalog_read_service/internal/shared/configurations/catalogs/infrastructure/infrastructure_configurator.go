package infrastructure

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/fxapp/contracts"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/pipelines"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/metrics"
	pipelines2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/metrics/mediatr/pipelines"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing"
	tracingpipelines "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing/mediatr/pipelines"

	"github.com/mehdihadeli/go-mediatr"
)

type InfrastructureConfigurator struct {
	contracts.Application
}

func NewInfrastructureConfigurator(
	app contracts.Application,
) *InfrastructureConfigurator {
	return &InfrastructureConfigurator{
		Application: app,
	}
}

func (ic *InfrastructureConfigurator) ConfigInfrastructures() {
	ic.ResolveFunc(
		func(logger logger.Logger, tracer tracing.AppTracer, metrics metrics.AppMetrics) error {
			err := mediatr.RegisterRequestPipelineBehaviors(
				pipelines.NewLoggingPipeline(logger),
				tracingpipelines.NewMediatorTracingPipeline(
					tracer,
					tracingpipelines.WithLogger(logger),
				),
				pipelines2.NewMediatorMetricsPipeline(
					metrics,
					pipelines2.WithLogger(logger),
				),
			)

			return err
		},
	)
}
