package otelTracer

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/otel/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

// ref:https://github.com/open-telemetry/opentelemetry-go-contrib/blob/df16f32df86b40077c9c90d06f33c4cdb6dd5afa/instrumentation/github.com/labstack/echo/otelecho/echo.go
// some changes in base code for handling 4xx error range to `ERROR` span state instead of `UNSET`
func Middleware(serviceName string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			request := c.Request()
			ctx := request.Context()

			//https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/semantic_conventions/http.md
			ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(request.Header))
			opts := []trace.SpanStartOption{
				trace.WithAttributes(semconv.NetAttributesFromHTTPRequest("tcp", request)...),
				trace.WithAttributes(semconv.EndUserAttributesFromHTTPRequest(request)...),
				trace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest(serviceName, c.Path(), request)...),
				trace.WithSpanKind(trace.SpanKindServer),
			}

			ctx, span := otel.Tracer("ehco").Start(ctx, fmt.Sprintf("%s process", c.Path()), opts...)
			defer span.End()

			//pass new ctx to next middleware
			c.SetRequest(request.WithContext(ctx))

			err := next(c)
			if err != nil {
				err = tracing.TraceErrFromSpan(span, err)
			}

			return err
		}
	}
}
