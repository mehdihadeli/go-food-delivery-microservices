package pipelines

import (
	"context"
	"fmt"
	"strings"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/constants/telemetry_attributes/app"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/constants/tracing/components"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing"
	attribute2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing/attribute"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing/utils"
	typeMapper "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/reflection/type_mappper"

	"emperror.dev/errors"
	"github.com/mehdihadeli/go-mediatr"
	"go.opentelemetry.io/otel/attribute"
)

type mediatorTracingPipeline struct {
	config *config
	tracer tracing.AppTracer
}

func NewMediatorTracingPipeline(
	appTracer tracing.AppTracer,
	opts ...Option,
) mediatr.PipelineBehavior {
	cfg := defaultConfig
	for _, opt := range opts {
		opt.apply(cfg)
	}

	return &mediatorTracingPipeline{
		config: cfg,
		tracer: appTracer,
	}
}

func (r *mediatorTracingPipeline) Handle(
	ctx context.Context,
	request interface{},
	next mediatr.RequestHandlerFunc,
) (interface{}, error) {
	requestName := typeMapper.GetSnakeTypeName(request)

	componentName := components.RequestHandler
	requestNameAttribute := app.RequestName
	requestAttribute := app.Request
	requestResultName := app.RequestResultName
	requestResult := app.RequestResult

	switch {
	case strings.Contains(typeMapper.GetPackageName(request), "command") || strings.Contains(typeMapper.GetPackageName(request), "commands"):
		componentName = components.CommandHandler
		requestNameAttribute = app.CommandName
		requestAttribute = app.Command
		requestResultName = app.CommandResultName
		requestResult = app.CommandResult
	case strings.Contains(typeMapper.GetPackageName(request), "query") || strings.Contains(typeMapper.GetPackageName(request), "queries"):
		componentName = components.QueryHandler
		requestNameAttribute = app.QueryName
		requestAttribute = app.Query
		requestResultName = app.QueryResultName
		requestResult = app.QueryResult
	case strings.Contains(typeMapper.GetPackageName(request), "event") || strings.Contains(typeMapper.GetPackageName(request), "events"):
		componentName = components.EventHandler
		requestNameAttribute = app.EventName
		requestAttribute = app.Event
		requestResultName = app.EventResultName
		requestResult = app.EventResult
	}

	operationName := fmt.Sprintf("%s_handler", requestName)
	spanName := fmt.Sprintf(
		"%s.%s/%s",
		componentName,
		operationName,
		requestName,
	) // by convention

	// https://golang.org/pkg/context/#Context
	newCtx, span := r.tracer.Start(ctx, spanName)

	defer span.End()

	span.SetAttributes(
		attribute.String(requestNameAttribute, requestName),
		attribute2.Object(requestAttribute, request),
	)

	response, err := next(newCtx)

	responseName := typeMapper.GetSnakeTypeName(response)
	span.SetAttributes(
		attribute.String(requestResultName, responseName),
		attribute2.Object(requestResult, response),
	)

	err = utils.TraceStatusFromSpan(
		span,
		errors.WrapIf(
			err,
			fmt.Sprintf("Request '%s' failed.", requestName),
		),
	)

	return response, err
}
