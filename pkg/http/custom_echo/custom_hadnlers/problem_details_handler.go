package customHadnlers

import (
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/http_errors"
	"net/http"
)
import "schneider.vip/problem"

func ProblemHandler(err error, c echo.Context) {
	if prb, ok := err.(problem.Problem); ok {
		if !c.Response().Committed {
			if _, err := prb.WriteTo(c.Response()); err != nil {
				c.Logger().Error(err)
			}
		}
	}
	if problemDetail, ok := err.(httpErrors.ProblemDetailErr); ok {
		if !c.Response().Committed {

			prb := problem.Of(problemDetail.GetStatus()).Append(problem.Detail(problemDetail.GetDetailError())).Append(problem.Title(problemDetail.GetTitle()))
			if _, err := prb.WriteTo(c.Response()); err != nil {
				c.Logger().Error(err)
			}
		}
	} else if echoErr, ok := err.(*echo.HTTPError); ok {
		if !c.Response().Committed {
			prb := problem.Of(echoErr.Code).Append(problem.Detail(echoErr.Message.(string)))
			if _, err := prb.WriteTo(c.Response()); err != nil {
				c.Logger().Error(err)
			}
		}
	} else {
		if !c.Response().Committed {
			prb := problem.Of(http.StatusInternalServerError).Append(problem.Detail(err.Error()))
			if _, err := prb.WriteTo(c.Response()); err != nil {
				c.Logger().Error(err)
			}

			// or
			//c.JSON(c.Response().Status, err)
		}
	}
}
