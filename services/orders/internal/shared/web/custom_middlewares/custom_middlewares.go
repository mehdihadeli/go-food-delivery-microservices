package cutomMiddlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/config"
)

type CustomMiddlewares interface {
	RequestLoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc
	RequestMetricsMiddleware(next echo.HandlerFunc) echo.HandlerFunc
}

type customMiddlewares struct {
	log         logger.Logger
	cfg         *config.Config
	metricsFunc MetricsFunc
}

func NewCustomMiddlewares(log logger.Logger, cfg *config.Config, metricsFunc MetricsFunc) *customMiddlewares {
	return &customMiddlewares{log: log, cfg: cfg, metricsFunc: metricsFunc}
}
