package customErrors

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/http_errors/contracts"
	"github.com/pkg/errors"
)

func NewApiError(message string, code int) *apiError {
	ae := &apiError{
		customError: NewCustomError(nil, code, message),
	}

	return ae
}

func NewApiErrorWrap(err error, code int, message string) *apiError {
	ae := &apiError{
		customError: NewCustomError(err, code, message),
	}

	return ae
}

type apiError struct {
	*customError
}

type ApiError interface {
	CustomError
	contracts.StackError
	IsApiError() bool
}

func (a *apiError) IsApiError() bool {
	return true
}

func (a *apiError) WithStack() error {
	// with this we use `Cause`, `Unwrap` method of new stack error but this struct `Cause`, `Unwrap` will call with next `Unwrap` on this object
	// Format this error (stackErr) with sprintf, First write Causer of error and then will write call stack for this point of code
	return errors.WithStack(a)
}

func IsApiError(err error) bool {
	var apiError ApiError

	//us, ok := errors.Cause(err).(applicationError)
	if errors.As(err, &apiError) {
		return apiError.IsApiError()
	}

	return false
}
