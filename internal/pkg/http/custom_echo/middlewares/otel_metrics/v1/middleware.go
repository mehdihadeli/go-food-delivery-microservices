package metricecho

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
)

// HTTPMetrics is an echo middleware to add metrics to rec for each HTTP request.
// If recorder config is nil, the middleware will use a recorder with default configuration.
func HTTPMetrics(cfg *HTTPRecorderConfig) echo.MiddlewareFunc {
	if cfg == nil {
		cfg = &HTTPCfg
	}

	rec := NewHTTPRecorder(*cfg, nil)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			values := HTTPLabels{
				Method: c.Request().Method,
				Path:   c.Path(),
			}

			rec.AddInFlightRequest(context.Background(), values)

			start := time.Now()

			defer func() {
				elapsed := time.Since(start)

				if err != nil {
					c.Error(err)
					// don't return the error so that it's not handled again
					err = nil
				}

				values.Code = c.Response().Status

				rec.AddRequestToTotal(context.Background(), values)
				rec.AddRequestDuration(context.Background(), elapsed, values)
				rec.RemInFlightRequest(context.Background(), values)
			}()

			return next(c)
		}
	}
}
