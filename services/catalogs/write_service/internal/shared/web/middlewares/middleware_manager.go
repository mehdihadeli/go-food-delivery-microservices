package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/config"
)

type MiddlewareManager interface {
	RequestLoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc
	RequestMetricsMiddleware(next echo.HandlerFunc) echo.HandlerFunc
}

type middlewareManager struct {
	log       logger.Logger
	cfg       *config.Config
	metricsCb MiddlewareMetricsCb
}

func NewMiddlewareManager(log logger.Logger, cfg *config.Config, metricsCb MiddlewareMetricsCb) *middlewareManager {
	return &middlewareManager{log: log, cfg: cfg, metricsCb: metricsCb}
}
