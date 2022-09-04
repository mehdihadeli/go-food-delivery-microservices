package problemDetails

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/constants"
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
		Title:      constants.ErrInternalServerErrorTitle,
		Detail:     detail,
		Status:     http.StatusInternalServerError,
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
