package customHadnlers

import (
	"github.com/labstack/echo/v4"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/problemDetails"
	defaultLogger "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/default_logger"
	errorUtils "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils/error_utils"
)

func ProblemHandler(err error, c echo.Context) {
	prb := problemDetails.ParseError(err)

	if prb != nil {
		if !c.Response().Committed {
			if _, err := problemDetails.WriteTo(prb, c.Response()); err != nil {
				defaultLogger.Logger.Error(err)
			}
		}
	} else {
		if !c.Response().Committed {
			prb := problemDetails.NewInternalServerProblemDetail(err.Error(), errorUtils.ErrorsWithStack(err))
			if _, err := problemDetails.WriteTo(prb, c.Response()); err != nil {
				defaultLogger.Logger.Error(err)
			}
		}
	}
}
