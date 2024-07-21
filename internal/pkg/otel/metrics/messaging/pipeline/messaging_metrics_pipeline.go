package pipelines

import (
	"context"
	"fmt"
	"time"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/messaging/pipeline"
	types2 "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/messaging/types"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/otel/constants/telemetrytags"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/otel/metrics"
	attribute2 "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/otel/tracing/attribute"

	"github.com/iancoleman/strcase"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type messagingMetricsPipeline struct {
	config *config
	meter  metrics.AppMetrics
}

func NewMessagingMetricsPipeline(
	appMetrics metrics.AppMetrics,
	opts ...Option,
) pipeline.ConsumerPipeline {
	cfg := defaultConfig
	for _, opt := range opts {
		opt.apply(cfg)
	}

	return &messagingMetricsPipeline{
		config: cfg,
		meter:  appMetrics,
	}
}

func (m *messagingMetricsPipeline) Handle(
	ctx context.Context,
	consumerContext types2.MessageConsumeContext,
	next pipeline.ConsumerHandlerFunc,
) error {
	message := consumerContext.Message()
	messageTypeName := message.GetMessageTypeName()
	snakeTypeName := strcase.ToSnake(messageTypeName)

	successRequestsCounter, err := m.meter.Int64Counter(
		fmt.Sprintf("%s.success_total", snakeTypeName),
		metric.WithUnit("count"),
		metric.WithDescription(
			fmt.Sprintf(
				"Measures the number of success '%s' (%s)",
				snakeTypeName,
				messageTypeName,
			),
		),
	)
	if err != nil {
		return err
	}

	failedRequestsCounter, err := m.meter.Int64Counter(
		fmt.Sprintf("%s.failed_total", snakeTypeName),
		metric.WithUnit("count"),
		metric.WithDescription(
			fmt.Sprintf(
				"Measures the number of failed '%s' (%s)",
				snakeTypeName,
				messageTypeName,
			),
		),
	)
	if err != nil {
		return err
	}

	totalRequestsCounter, err := m.meter.Int64Counter(
		fmt.Sprintf("%s.total", snakeTypeName),
		metric.WithUnit("count"),
		metric.WithDescription(
			fmt.Sprintf(
				"Measures the total number of '%s' (%s)",
				snakeTypeName,
				messageTypeName,
			),
		),
	)
	if err != nil {
		return err
	}

	durationValueRecorder, err := m.meter.Int64Histogram(
		fmt.Sprintf("%s.duration", snakeTypeName),
		metric.WithUnit("ms"),
		metric.WithDescription(
			fmt.Sprintf(
				"Measures the duration of '%s' (%s)",
				snakeTypeName,
				messageTypeName,
			),
		),
	)
	if err != nil {
		return err
	}

	// Start recording the duration
	startTime := time.Now()

	err = next(ctx)

	// Calculate the duration
	duration := time.Since(startTime).Milliseconds()

	opt := metric.WithAttributes(
		attribute.String(telemetrytags.App.MessageType, messageTypeName),
		attribute.String(telemetrytags.App.MessageName, snakeTypeName),
		attribute2.Object(telemetrytags.App.Message, message),
	)

	// Record metrics
	totalRequestsCounter.Add(ctx, 1, opt)

	if err == nil {
		successRequestsCounter.Add(ctx, 1, opt)
	} else {
		failedRequestsCounter.Add(ctx, 1, opt)
	}

	durationValueRecorder.Record(ctx, duration, opt)

	return nil
}
