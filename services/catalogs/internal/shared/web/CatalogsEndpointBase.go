package web

import (
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/shared"
	"github.com/opentracing/opentracing-go"
)

type CatalogsEndpointBase struct {
	Echo      *echo.Echo
	Log       logger.Logger
	Cfg       *config.Config
	Mediator  *mediatr.Mediator
	Validator *validator.Validate
	Metrics   *shared.CatalogsServiceMetrics
}

func (h *CatalogsEndpointBase) TraceErr(span opentracing.Span, err error) {
	span.SetTag("error", true)
	span.LogKV("error_code", err.Error())
	h.Metrics.ErrorHttpRequests.Inc()
}
