package otelmetrics

// ref:https://github.com/open-telemetry/opentelemetry-go/blob/main/example/prometheus/main.go
// https://github.com/labstack/echo-contrib/blob/master/prometheus/prometheus.go
// https://github.com/worldline-go/tell/tree/main/metric/metricecho
// https://opentelemetry.io/docs/instrumentation/go/manual/#metrics

// https://opentelemetry.io/docs/specs/otel/metrics/semantic_conventions/http-metrics/

import (
	"time"

	"github.com/labstack/echo/v4"
)

// HTTPMetrics is a middleware for adding  otel metrics for a given request
// If recorder config is nil, the middleware will use a recorder with default configuration.
func HTTPMetrics(opts ...Option) echo.MiddlewareFunc {
	config := defualtConfig

	for _, opt := range opts {
		opt.apply(&config)
	}

	httpMetricsRecorder := NewHTTPMetricsRecorder(config)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if config.skipper(c) {
				return next(c)
			}

			request := c.Request()
			ctx := request.Context()

			values := HTTPLabels{
				Method: request.Method,
				Path:   c.Path(),
				Host:   request.URL.Host,
			}

			httpMetricsRecorder.AddInFlightRequest(ctx, values)

			start := time.Now()

			defer func() {
				elapsed := time.Since(start)

				values.Code = c.Response().Status

				httpMetricsRecorder.AddRequestToTotal(ctx, values)

				httpMetricsRecorder.AddRequestDuration(ctx, elapsed, values)

				httpMetricsRecorder.RemInFlightRequest(ctx, values)

				httpMetricsRecorder.AddRequestSize(ctx, request, values)

				httpMetricsRecorder.AddResponseSize(ctx, c.Response(), values)

				if err != nil {
					httpMetricsRecorder.AddRequestError(ctx, values)
				} else {
					httpMetricsRecorder.AddRequestSuccess(ctx, values)
				}
			}()

			err = next(c)
			if err != nil {
				// handle echo error in this middleware and raise echo errorhandler func and our custom error handler
				// when we call c.Error more than once, `c.Response().Committed` becomes true and response doesn't write to client again in our error handler
				// Error will update response status with occurred error object status code
				c.Error(err)
			}

			return err
		}
	}
}
