package oteltracing

// Ref: https://github.com/open-telemetry/opentelemetry-go-contrib/blob/main/instrumentation/github.com/labstack/echo/otelecho/echo.go
// Note: for consideration of 4xx status as error in traces, I customized original echo oteltrace middleware for handing my requirements

// https://opentelemetry.io/docs/specs/otel/trace/semantic_conventions/http/

import (
	"fmt"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/otel/tracing/utils"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/semconv/v1.20.0/httpconv"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// HttpTrace returns echo middleware which will trace incoming requests.
func HttpTrace(opts ...Option) echo.MiddlewareFunc {
	cfg := defualtConfig
	for _, opt := range opts {
		opt.apply(&cfg)
	}

	tracer := cfg.tracerProvider.Tracer(cfg.instrumentationName)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if cfg.skipper(c) {
				return next(c)
			}

			c.Set(cfg.instrumentationName, tracer)
			request := c.Request()
			// doesn't contain trace information and after completing trace on new ctx we should go back to our old savedCtx
			savedCtx := request.Context()

			defer func() {
				// we should go back to previous context in end of our operation because new context contains child spans, and if we don't set it back to previous context, after returning from this method all further parent spans becomes a new child for existing child span!
				request = request.WithContext(savedCtx)
				c.SetRequest(request)
			}()

			// create new ctx from existing savedCtx
			ctx := cfg.propagators.Extract(
				savedCtx,
				propagation.HeaderCarrier(request.Header),
			)

			// //https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/semantic_conventions/http.md
			// httpconv doesn't exist in semconv v1.21.0 we have to use v1.20.0 for that
			// https://github.com/open-telemetry/opentelemetry-go/pull/4362
			opts := []oteltrace.SpanStartOption{
				oteltrace.WithAttributes(
					httpconv.ServerRequest(cfg.serviceName, request)...),
				oteltrace.WithSpanKind(oteltrace.SpanKindServer),
			}

			if path := c.Path(); path != "" {
				rAttr := semconv.HTTPRoute(path)
				opts = append(opts, oteltrace.WithAttributes(rAttr))
			}

			spanName := c.Path()
			if spanName == "" {
				spanName = fmt.Sprintf(
					"HTTP %s route not found",
					request.Method,
				)
			}

			ctx, span := tracer.Start(ctx, spanName, opts...)
			defer span.End()

			// add the new context into the request, because new ctx contains our created span and we want all inner spans in next middlewares become child span of this span
			// pass the span through the request context
			c.SetRequest(request.WithContext(ctx))

			// serve the request to the next middleware
			err := next(c)
			if err != nil {
				// handle echo error in this middleware and raise echo errorhandler func and our custom error handler
				// when we call c.Error more than once, `c.Response().Committed` becomes true and response doesn't write to client again in our error handler
				// Error will update response status with occurred error object status code
				c.Error(err)
			}

			status := c.Response().Status
			err = utils.HttpTraceStatusFromSpanWithCode(span, err, status)

			return err
		}
	}
}
