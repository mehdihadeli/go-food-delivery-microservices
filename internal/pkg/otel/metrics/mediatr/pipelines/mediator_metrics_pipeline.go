package pipelines

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/constants/telemetry_attributes/app"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/metrics"
	attribute2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing/attribute"
	typeMapper "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/reflection/type_mappper"

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
	requestName := typeMapper.GetSnakeTypeName(request)

	requestNameAttribute := app.RequestName
	requestAttribute := app.Request
	requestResultName := app.RequestResultName
	requestResult := app.RequestResult
	requestType := "request"

	switch {
	case strings.Contains(typeMapper.GetPackageName(request), "command") || strings.Contains(typeMapper.GetPackageName(request), "commands"):
		requestNameAttribute = app.CommandName
		requestAttribute = app.Command
		requestResultName = app.CommandResultName
		requestResult = app.CommandResult
		requestType = "command"
	case strings.Contains(typeMapper.GetPackageName(request), "query") || strings.Contains(typeMapper.GetPackageName(request), "queries"):
		requestNameAttribute = app.QueryName
		requestAttribute = app.Query
		requestResultName = app.QueryResultName
		requestResult = app.QueryResult
		requestType = "query"
	case strings.Contains(typeMapper.GetPackageName(request), "event") || strings.Contains(typeMapper.GetPackageName(request), "events"):
		requestNameAttribute = app.EventName
		requestAttribute = app.Event
		requestResultName = app.EventResultName
		requestResult = app.EventResult
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

	responseName := typeMapper.GetSnakeTypeName(response)
	opt := metric.WithAttributes(
		attribute.String(requestNameAttribute, requestName),
		attribute2.Object(requestAttribute, request),
		attribute.String(requestResultName, responseName),
		attribute2.Object(requestResult, response),
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
