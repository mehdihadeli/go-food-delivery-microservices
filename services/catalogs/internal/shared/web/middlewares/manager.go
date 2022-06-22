package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/config"
	"strings"
	"time"
)

type MiddlewareManager interface {
	RequestLoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc
}

type MiddlewareMetricsCb func(err error)

type middlewareManager struct {
	log       logger.Logger
	cfg       *config.Config
	metricsCb MiddlewareMetricsCb
}

func NewMiddlewareManager(log logger.Logger, cfg *config.Config, metricsCb MiddlewareMetricsCb) *middlewareManager {
	return &middlewareManager{log: log, cfg: cfg, metricsCb: metricsCb}
}

func (mw *middlewareManager) RequestLoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {

		start := time.Now()
		err := next(ctx)

		req := ctx.Request()
		res := ctx.Response()
		status := res.Status
		size := res.Size
		s := time.Since(start)

		if !mw.checkIgnoredURI(ctx.Request().RequestURI, mw.cfg.Http.IgnoreLogUrls) {
			mw.log.HttpMiddlewareAccessLogger(req.Method, req.URL.String(), status, size, s)
		}

		mw.metricsCb(err)
		return err
	}
}

func (mw *middlewareManager) checkIgnoredURI(requestURI string, uriList []string) bool {
	for _, s := range uriList {
		if strings.Contains(requestURI, s) {
			return true
		}
	}
	return false
}
