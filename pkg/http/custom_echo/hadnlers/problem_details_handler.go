package customHadnlers

import (
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/problemDetails"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger/defaultLogger"
	errorUtils "github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils/error_utils"
)

func ProblemHandler(err error, c echo.Context) {
	prb := problemDetails.ParseError(err)

	if prb != nil {
		if !c.Response().Committed {
			if _, err := prb.WriteTo(c.Response()); err != nil {
				defaultLogger.Logger.Error(err)
			}
		}
	} else {
		if !c.Response().Committed {
			prb := problemDetails.NewInternalServerProblemDetail(err.Error(), errorUtils.ErrorsWithStack(err))
			if _, err := prb.WriteTo(c.Response()); err != nil {
				defaultLogger.Logger.Error(err)
			}
		}
	}
}
