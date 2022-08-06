package cutomMiddlewares

import (
	"github.com/labstack/echo/v4"
)

type MetricsFunc func(err error)

func (mw *customMiddlewares) RequestMetricsMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		err := next(ctx)
		mw.metricsFunc(err)

		return err
	}
}
