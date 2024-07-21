package tracing

import (
	"context"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/otel/tracing/utils"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type AppTracer interface {
	trace.Tracer
}

type appTracer struct {
	trace.Tracer
}

func (c *appTracer) Start(
	ctx context.Context,
	spanName string,
	opts ...trace.SpanStartOption,
) (context.Context, trace.Span) {
	parentSpan := trace.SpanFromContext(ctx)
	if parentSpan != nil {
		utils.ContextWithParentSpan(ctx, parentSpan)
	}

	return c.Tracer.Start(ctx, spanName, opts...)
}

func NewAppTracer(name string, options ...trace.TracerOption) AppTracer {
	// without registering `NewOtelTracing` it uses global empty (NoopTracer) TraceProvider but after using `NewOtelTracing`, global TraceProvider will be replaced
	tracer := otel.Tracer(name, options...)
	return &appTracer{Tracer: tracer}
}
