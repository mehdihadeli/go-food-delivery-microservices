package problem_details

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/constants"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/http_errors"
	"github.com/pkg/errors"
	"net/http"
	"strings"
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
	if problemDetail, ok := err.(http_errors.ProblemDetailErr); ok {
		if !c.Response().Committed {
			prb := problem.New(problem.Title(problemDetail.GetTitle()), problem.Status(problemDetail.GetStatus()), problem.Detail(problemDetail.GetDetail()))
			if _, err := prb.WriteTo(c.Response()); err != nil {
				c.Logger().Error(err)
			}
		}
	} else {
		problemDetail := MapErrors(err)
		if !c.Response().Committed {
			prb := problem.New(problem.Title(problemDetail.GetTitle()), problem.Status(problemDetail.GetStatus()), problem.Detail(problemDetail.GetDetail()))
			if _, err := prb.WriteTo(c.Response()); err != nil {
				c.Logger().Error(err)
			}
		}
	}

	// Take required information from error and context and send it to a service like New Relic
	fmt.Println(c.Path(), c.QueryParams(), err.Error())
}

// MapErrors  map of error string messages returns ProblemDetailErr
func MapErrors(err error) http_errors.ProblemDetailErr {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return http_errors.NewProblemDetailError(http.StatusNotFound, constants.ErrNotFound, err.Error())
	case errors.Is(err, context.DeadlineExceeded):
		return http_errors.NewProblemDetailError(http.StatusRequestTimeout, constants.ErrRequestTimeout, err.Error())
	case strings.Contains(strings.ToLower(err.Error()), constants.SQLState):
		return parseSqlErrors(err)
	case strings.Contains(strings.ToLower(err.Error()), "field validation"):
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return http_errors.NewProblemDetailError(http.StatusBadRequest, constants.ErrBadRequest, validationErrors.Error())
		}
		return parseValidatorError(err)
	case strings.Contains(strings.ToLower(err.Error()), "required header"):
		return http_errors.NewProblemDetailError(http.StatusBadRequest, constants.ErrBadRequest, err.Error())
	case strings.Contains(strings.ToLower(err.Error()), constants.Base64):
		return http_errors.NewProblemDetailError(http.StatusBadRequest, constants.ErrBadRequest, err.Error())
	case strings.Contains(strings.ToLower(err.Error()), constants.Unmarshal):
		return http_errors.NewProblemDetailError(http.StatusBadRequest, constants.ErrBadRequest, err.Error())
	case strings.Contains(strings.ToLower(err.Error()), constants.Uuid):
		return http_errors.NewProblemDetailError(http.StatusBadRequest, constants.ErrBadRequest, err.Error())
	case strings.Contains(strings.ToLower(err.Error()), constants.Cookie):
		return http_errors.NewProblemDetailError(http.StatusUnauthorized, constants.ErrUnauthorized, err.Error())
	case strings.Contains(strings.ToLower(err.Error()), constants.Token):
		return http_errors.NewProblemDetailError(http.StatusUnauthorized, constants.ErrUnauthorized, err.Error())
	case strings.Contains(strings.ToLower(err.Error()), constants.Bcrypt):
		return http_errors.NewProblemDetailError(http.StatusBadRequest, constants.ErrBadRequest, err.Error())
	case strings.Contains(strings.ToLower(err.Error()), "no documents in result"):
		return http_errors.NewProblemDetailError(http.StatusNotFound, constants.ErrNotFound, err.Error())
	default:
		if problemDetailErr, ok := err.(http_errors.ProblemDetail); ok {
			return problemDetailErr
		}
		return http_errors.NewProblemDetailError(http.StatusInternalServerError, constants.ErrInternalServerError, errors.Cause(err).Error())
	}
}

func parseSqlErrors(err error) http_errors.ProblemDetailErr {
	return http_errors.NewProblemDetailError(http.StatusBadRequest, constants.ErrBadRequest, err.Error())
}

func parseValidatorError(err error) http_errors.ProblemDetailErr {
	if strings.Contains(err.Error(), "Password") {
		return http_errors.NewProblemDetailError(http.StatusBadRequest, constants.ErrInvalidPassword, err.Error())
	}

	if strings.Contains(err.Error(), "Email") {
		return http_errors.NewProblemDetailError(http.StatusBadRequest, constants.ErrInvalidEmail, err.Error())
	}

	return http_errors.NewProblemDetailError(http.StatusBadRequest, constants.ErrInvalidField, err.Error())
}
