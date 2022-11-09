package otelMetrics

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
)

var (
	successCounter syncfloat64.Counter
	errorCounter   syncfloat64.Counter
)

// Middleware adds request status metrics to the otel
// ref: https://github.com/open-telemetry/opentelemetry-go/blob/main/example/prometheus/main.go
func Middleware(meter metric.Meter, serviceName string) echo.MiddlewareFunc {
	errorCounter, _ = meter.SyncFloat64().Counter(fmt.Sprintf("%s_error_http_requests_total", serviceName), instrument.WithDescription("The total number of error http requests"))
	successCounter, _ = meter.SyncFloat64().Counter(fmt.Sprintf("%s_success_http_requests_total", serviceName), instrument.WithDescription("The total number of success http requests"))

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)

			request := c.Request()
			ctx := request.Context()

			attrs := []attribute.KeyValue{
				attribute.Key("MetricsType").String("Http"),
			}

			if err != nil {
				c.Error(err)
				if err != nil {
					return err
				}
				errorCounter.Add(ctx, 1, attrs...)
			} else {
				if err != nil {
					return err
				}
				successCounter.Add(ctx, 1, attrs...)
			}

			// update request context
			c.SetRequest(request.WithContext(ctx))

			return err
		}
	}
}
