package tracing

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type CustomTracer interface {
	trace.Tracer
}

type customTracer struct {
	trace.Tracer
}

func (c *customTracer) Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	parentSpan := trace.SpanFromContext(ctx)
	if parentSpan != nil {
		ContextWithParentSpan(ctx, parentSpan)
	}

	return c.Tracer.Start(ctx, spanName, opts...)
}

func NewCustomTracer(name string, options ...trace.TracerOption) CustomTracer {
	tracer := otel.Tracer(name, options...)
	return &customTracer{Tracer: tracer}
}
