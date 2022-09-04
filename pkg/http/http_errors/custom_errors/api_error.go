package customErrors

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/contracts"
	"github.com/pkg/errors"
)

func NewApiError(message string, code int) error {
	ae := &apiError{
		WithStack: NewCustomErrorStack(nil, code, message),
	}

	return ae
}

func NewApiErrorWrap(err error, code int, message string) error {
	ae := &apiError{
		WithStack: NewCustomErrorStack(err, code, message),
	}

	return ae
}

type apiError struct {
	contracts.WithStack
}

type ApiError interface {
	contracts.WithStack
	GetCustomError() CustomError
	IsApiError() bool
}

func (a *apiError) IsApiError() bool {
	return true
}

func (a *apiError) GetCustomError() CustomError {
	return GetCustomError(a)
}

func IsApiError(err error) bool {
	a, ok := err.(ApiError)
	if ok && a.IsApiError() {
		return true
	}

	var apiError ApiError
	//us, ok := grpc_errors.Cause(err).(ApiError)
	if errors.As(err, &apiError) {
		return apiError.IsApiError()
	}

	return false
}
