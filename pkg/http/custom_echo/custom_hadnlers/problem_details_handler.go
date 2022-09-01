package customHadnlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core"
	httpErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/http_errors/problemDetails"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger/defaultLogger"
)

func ProblemHandler(err error, c echo.Context) {
	prb := problemDetails.ParseError(err)

	if prb != nil {
		if !c.Response().Committed {
			if _, err := prb.WriteTo(c.Response()); err != nil {
				defaultLogger.Logger.Error(err)
			}
			defaultLogger.Logger.Error(prb.Error())
			if core.IsDevelopment() {
				stackTrace := prb.GetStackTrace()
				fmt.Println(stackTrace)
			}
		}
	} else {
		if !c.Response().Committed {
			prb := problemDetails.NewInternalServerProblemDetail(err.Error(), httpErrors.ErrorsWithStack(err))
			if _, err := prb.WriteTo(c.Response()); err != nil {
				defaultLogger.Logger.Error(err)
			}
			defaultLogger.Logger.Error(prb.Error())
			if core.IsDevelopment() {
				stackTrace := prb.GetStackTrace()
				fmt.Println(stackTrace)
			}
		}
	}
}
