package problemDetails

import (
	"context"
	"database/sql"
	"github.com/go-playground/validator"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/constants"
	httpErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http_errors"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http_errors/custom_errors"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

func NewValidationProblemDetail(detail string, stackTrace string) ProblemDetailErr {
	validationError :=
		&problemDetail{
			Title:      constants.ErrBadRequestTitle,
			Detail:     detail,
			Status:     http.StatusBadRequest,
			Type:       getDefaultType(http.StatusBadRequest),
			Timestamp:  time.Now(),
			StackTrace: stackTrace,
		}

	return validationError
}

func NewConflictProblemDetail(detail string, stackTrace string) ProblemDetailErr {
	return &problemDetail{
		Title:      constants.ErrConflictTitle,
		Detail:     detail,
		Status:     http.StatusConflict,
		Type:       getDefaultType(http.StatusConflict),
		Timestamp:  time.Now(),
		StackTrace: stackTrace,
	}
}

func NewBadRequestProblemDetail(detail string, stackTrace string) ProblemDetailErr {
	return &problemDetail{
		Title:      constants.ErrBadRequestTitle,
		Detail:     detail,
		Status:     http.StatusBadRequest,
		Type:       getDefaultType(http.StatusBadRequest),
		Timestamp:  time.Now(),
		StackTrace: stackTrace,
	}
}

func NewWrongCredentialsProblemDetail(detail string, stackTrace string) ProblemDetailErr {
	return &problemDetail{
		Title:      constants.ErrWrongCredentialsTitle,
		Detail:     detail,
		Status:     http.StatusUnauthorized,
		Type:       getDefaultType(http.StatusUnauthorized),
		Timestamp:  time.Now(),
		StackTrace: stackTrace,
	}
}

func NewNotFoundErrorProblemDetail(detail string, stackTrace string) ProblemDetailErr {
	return &problemDetail{
		Title:      constants.ErrNotFoundTitle,
		Detail:     detail,
		Status:     http.StatusNotFound,
		Type:       getDefaultType(http.StatusNotFound),
		Timestamp:  time.Now(),
		StackTrace: stackTrace,
	}
}

func NewUnAuthorizedErrorProblemDetail(detail string, stackTrace string) ProblemDetailErr {
	return &problemDetail{
		Title:      constants.ErrUnauthorizedTitle,
		Detail:     detail,
		Status:     http.StatusUnauthorized,
		Type:       getDefaultType(http.StatusUnauthorized),
		Timestamp:  time.Now(),
		StackTrace: stackTrace,
	}
}

func NewForbiddenProblemDetail(detail string, stackTrace string) ProblemDetailErr {
	return &problemDetail{
		Title:      constants.ErrForbiddenTitle,
		Detail:     detail,
		Status:     http.StatusForbidden,
		Type:       getDefaultType(http.StatusForbidden),
		Timestamp:  time.Now(),
		StackTrace: stackTrace,
	}
}

func NewInternalServerProblemDetail(detail string, stackTrace string) ProblemDetailErr {
	return &problemDetail{
		Title:  constants.ErrInternalServerErrorTitle,
		Detail: detail, Status: http.StatusInternalServerError,
		Type:       getDefaultType(http.StatusInternalServerError),
		Timestamp:  time.Now(),
		StackTrace: stackTrace,
	}
}

func NewDomainProblemDetail(status int, detail string, stackTrace string) ProblemDetailErr {
	return &problemDetail{
		Title:      constants.ErrDomainTitle,
		Detail:     detail,
		Status:     status,
		Type:       getDefaultType(http.StatusBadRequest),
		Timestamp:  time.Now(),
		StackTrace: stackTrace,
	}
}

func NewApplicationProblemDetail(status int, detail string, stackTrace string) ProblemDetailErr {
	return &problemDetail{
		Title:      constants.ErrApplicationTitle,
		Detail:     detail,
		Status:     status,
		Type:       getDefaultType(status),
		Timestamp:  time.Now(),
		StackTrace: stackTrace,
	}
}

func NewApiProblemDetail(status int, detail string, stackTrace string) ProblemDetailErr {
	return &problemDetail{
		Title:      constants.ErrApiTitle,
		Detail:     detail,
		Status:     status,
		Type:       getDefaultType(status),
		Timestamp:  time.Now(),
		StackTrace: stackTrace,
	}
}

func ParseError(err error) ProblemDetailErr {
	stackTrace := httpErrors.ErrorsWithStack(err)
	customErr := customErrors.GetCustomError(err)
	var validatorErr validator.ValidationErrors

	if err != nil {
		switch {
		case customErrors.IsDomainError(err):
			return NewDomainProblemDetail(customErr.Status(), customErr.Error(), stackTrace)
		case customErrors.IsApplicationError(err):
			return NewApplicationProblemDetail(customErr.Status(), customErr.Error(), stackTrace)
		case customErrors.IsBadRequestError(err):
			return NewBadRequestProblemDetail(customErr.Error(), stackTrace)
		case customErrors.IsNotFoundError(err):
			return NewNotFoundErrorProblemDetail(customErr.Error(), stackTrace)
		case customErrors.IsValidationError(err):
			return NewValidationProblemDetail(customErr.Error(), stackTrace)
		case customErrors.IsUnAuthorizedError(err):
			return NewUnAuthorizedErrorProblemDetail(customErr.Error(), stackTrace)
		case customErrors.IsForbiddenError(err):
			return NewForbiddenProblemDetail(customErr.Error(), stackTrace)
		case customErrors.IsConflictError(err):
			return NewConflictProblemDetail(customErr.Error(), stackTrace)
		case customErrors.IsInternalServerError(err):
			return NewInternalServerProblemDetail(customErr.Error(), stackTrace)
		case customErrors.IsApiError(err):
			return NewApiProblemDetail(customErr.Status(), customErr.Error(), stackTrace)
		case customErrors.IsCustomError(err):
			return NewProblemDetailFromCodeAndDetail(customErr.Status(), customErr.Error(), stackTrace)
		case customErrors.IsUnMarshalingError(err):
			return NewInternalServerProblemDetail(err.Error(), stackTrace)
		case customErrors.IsMarshalingError(err):
			return NewInternalServerProblemDetail(err.Error(), stackTrace)
		case errors.Is(err, sql.ErrNoRows):
			return NewNotFoundErrorProblemDetail(err.Error(), stackTrace)
		case errors.Is(err, context.DeadlineExceeded):
			return NewProblemDetail(http.StatusRequestTimeout, constants.ErrRequestTimeoutTitle, err.Error(), stackTrace)
		case errors.As(err, &validatorErr):
			return NewValidationProblemDetail(validatorErr.Error(), stackTrace)
		default:
			return NewInternalServerProblemDetail(err.Error(), stackTrace)
		}
	}

	return nil
}
