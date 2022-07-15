package middlewares

import "github.com/labstack/echo/v4"

type MiddlewareMetricsCb func(err error)

func (mw *middlewareManager) RequestMetricsMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		err := next(ctx)
		mw.metricsCb(err)

		return err
	}
}
