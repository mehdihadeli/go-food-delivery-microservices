package tracing

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/otel/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

//https://opentelemetry.io/docs/reference/specification/trace/semantic_conventions/http/
//https://opentelemetry.io/docs/instrumentation/go/manual/#semantic-attributes

var HttpTracer trace.Tracer

func init() {
	HttpTracer = tracing.NewCustomTracer("github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/custom_echo") //instrumentation name
}

// SpanFromContext get current context span from existing echo otel trace middleware instrument
func SpanFromContext(c echo.Context) (span trace.Span) {
	ctx := c.Request().Context()
	span = trace.SpanFromContext(ctx)

	return span
}

// StartHttpTraceSpan uses when echo otel middleware is off and create a span on 'http-echo' tracer
func StartHttpTraceSpan(c echo.Context, operationName string) (ctx context.Context, span trace.Span, deferSpan func()) {
	request := c.Request()
	ctx = request.Context()

	//https://github.com/open-telemetry/opentelemetry-go-contrib/blob/main/instrumentation/github.com/labstack/echo/otelecho/echo.go
	//https://lightstep.com/blog/opentelemetry-go-all-you-need-to-know
	pro := otel.GetTextMapPropagator()
	ctx = pro.Extract(ctx, propagation.HeaderCarrier(c.Request().Header))

	opts := []trace.SpanStartOption{
		trace.WithAttributes(semconv.NetAttributesFromHTTPRequest("tcp", request)...),
		trace.WithAttributes(semconv.EndUserAttributesFromHTTPRequest(request)...),
		trace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest(c.Request().Host, c.Path(), request)...),
		trace.WithSpanKind(trace.SpanKindServer),
	}

	ctx, span = HttpTracer.Start(ctx, operationName, opts...)

	// pass the span through the request ctx
	c.SetRequest(request.WithContext(ctx))

	return ctx, span, func() {
		attrs := semconv.HTTPAttributesFromHTTPStatusCode(c.Response().Status)
		spanStatus, spanMessage := semconv.SpanStatusFromHTTPStatusCodeAndSpanKind(c.Response().Status, trace.SpanKindServer)
		span.SetAttributes(attrs...)
		span.SetStatus(spanStatus, spanMessage)

		span.End()
	}
}
