package otelMetrics

// ref:https://github.com/open-telemetry/opentelemetry-go/blob/main/example/prometheus/main.go

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
	api "go.opentelemetry.io/otel/metric"
)

var (
	successCounter api.Float64Counter
	errorCounter   api.Float64Counter
)

// Middleware adds request status metrics to the otel
// ref: https://github.com/open-telemetry/opentelemetry-go/blob/main/example/prometheus/main.go
func Middleware(meter api.Meter, serviceName string) echo.MiddlewareFunc {
	errorCounter, _ := meter.Float64Counter(
		fmt.Sprintf("%s_error_http_requests_total", serviceName),
		api.WithDescription("The total number of error http requests"),
	)
	successCounter, _ = meter.Float64Counter(
		fmt.Sprintf("%s_success_http_requests_total", serviceName),
		api.WithDescription("The total number of success http requests"),
	)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			request := c.Request()
			ctx := request.Context()

			attrs := api.WithAttributes(
				attribute.Key("MetricsType").String("Http"),
			)

			if err != nil {
				errorCounter.Add(ctx, 1, attrs)
			} else {
				successCounter.Add(ctx, 1, attrs)
			}

			// update request context
			c.SetRequest(request.WithContext(ctx))

			return err
		}
	}
}
