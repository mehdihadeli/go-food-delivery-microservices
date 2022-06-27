package problem_details

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/http_errors"
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
	if restE, ok := err.(http_errors.RestError); ok {
		if !c.Response().Committed {
			prb := problem.New(problem.Title(restE.ErrError), problem.Status(restE.ErrStatus), problem.Detail(restE.ErrMessage.(string)))
			if _, err := prb.WriteTo(c.Response()); err != nil {
				c.Logger().Error(err)
			}
		}
	} else {
		c.Echo().DefaultHTTPErrorHandler(err, c)
	}

	// Take required information from error and context and send it to a service like New Relic
	fmt.Println(c.Path(), c.QueryParams(), err.Error())
}
