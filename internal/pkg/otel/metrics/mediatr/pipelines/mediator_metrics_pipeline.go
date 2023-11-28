package pipelines

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/constants/telemetrytags"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/metrics"
	customAttribute "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing/attribute"
	typemapper "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/reflection/typemapper"

	"github.com/mehdihadeli/go-mediatr"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type mediatorMetricsPipeline struct {
	config *config
	meter  metrics.AppMetrics
}

func NewMediatorMetricsPipeline(
	appMetrics metrics.AppMetrics,
	opts ...Option,
) mediatr.PipelineBehavior {
	cfg := defaultConfig
	for _, opt := range opts {
		opt.apply(cfg)
	}

	return &mediatorMetricsPipeline{
		config: cfg,
		meter:  appMetrics,
	}
}

func (r *mediatorMetricsPipeline) Handle(
	ctx context.Context,
	request interface{},
	next mediatr.RequestHandlerFunc,
) (interface{}, error) {
	requestName := typemapper.GetSnakeTypeName(request)

	requestNameTag := telemetrytags.App.RequestName
	requestTag := telemetrytags.App.Request
	requestResultNameTag := telemetrytags.App.RequestResultName
	requestResultTag := telemetrytags.App.RequestResult
	requestType := "request"

	switch {
	case strings.Contains(typemapper.GetPackageName(request), "command") || strings.Contains(typemapper.GetPackageName(request), "commands"):
		requestNameTag = telemetrytags.App.CommandName
		requestTag = telemetrytags.App.Command
		requestResultNameTag = telemetrytags.App.CommandResultName
		requestResultTag = telemetrytags.App.CommandResult
		requestType = "command"
	case strings.Contains(typemapper.GetPackageName(request), "query") || strings.Contains(typemapper.GetPackageName(request), "queries"):
		requestNameTag = telemetrytags.App.QueryName
		requestTag = telemetrytags.App.Query
		requestResultNameTag = telemetrytags.App.QueryResultName
		requestResultTag = telemetrytags.App.QueryResult
		requestType = "query"
	case strings.Contains(typemapper.GetPackageName(request), "event") || strings.Contains(typemapper.GetPackageName(request), "events"):
		requestNameTag = telemetrytags.App.EventName
		requestTag = telemetrytags.App.Event
		requestResultNameTag = telemetrytags.App.EventResultName
		requestResultTag = telemetrytags.App.EventResult
		requestType = "event"
	}

	successRequestsCounter, err := r.meter.Int64Counter(
		fmt.Sprintf("%s.success_total", requestName),
		metric.WithUnit("count"),
		metric.WithDescription(
			fmt.Sprintf(
				"Measures the number of success '%s' %s",
				requestName,
				requestType,
			),
		),
	)
	if err != nil {
		return nil, err
	}

	failedRequestsCounter, err := r.meter.Int64Counter(
		fmt.Sprintf("%s.failed_total", requestName),
		metric.WithUnit("count"),
		metric.WithDescription(
			fmt.Sprintf(
				"Measures the number of failed '%s' %s",
				requestName,
				requestType,
			),
		),
	)
	if err != nil {
		return nil, err
	}

	totalRequestsCounter, err := r.meter.Int64Counter(
		fmt.Sprintf("%s.total", requestName),
		metric.WithUnit("count"),
		metric.WithDescription(
			fmt.Sprintf(
				"Measures the total number of '%s' %s",
				requestName,
				requestType,
			),
		),
	)
	if err != nil {
		return nil, err
	}

	durationValueRecorder, err := r.meter.Int64Histogram(
		fmt.Sprintf("%s.duration", requestName),
		metric.WithUnit("ms"),
		metric.WithDescription(
			fmt.Sprintf(
				"Measures the duration of '%s' %s",
				requestName,
				requestType,
			),
		),
	)
	if err != nil {
		return nil, err
	}

	// Start recording the duration
	startTime := time.Now()

	response, err := next(ctx)

	// Calculate the duration
	duration := time.Since(startTime).Milliseconds()

	// response will be nil if we have an error
	responseName := typemapper.GetSnakeTypeName(response)

	opt := metric.WithAttributes(
		attribute.String(requestNameTag, requestName),
		customAttribute.Object(requestTag, request),
		attribute.String(requestResultNameTag, responseName),
		customAttribute.Object(requestResultTag, response),
	)

	// Record metrics
	totalRequestsCounter.Add(ctx, 1, opt)

	if err == nil {
		successRequestsCounter.Add(ctx, 1, opt)
	} else {
		failedRequestsCounter.Add(ctx, 1, opt)
	}

	durationValueRecorder.Record(ctx, duration, opt)

	return response, err
}
