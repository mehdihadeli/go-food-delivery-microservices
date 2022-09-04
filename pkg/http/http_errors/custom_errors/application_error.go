package customErrors

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/contracts"
	"github.com/pkg/errors"
	"net/http"
)

func NewApplicationError(message string) error {
	ae := &applicationError{
		WithStack: NewCustomErrorStack(nil, http.StatusInternalServerError, message),
	}

	return ae
}

func NewApplicationErrorWithCode(message string, code int) error {
	ae := &applicationError{
		WithStack: NewCustomErrorStack(nil, code, message),
	}

	return ae
}

func NewApplicationErrorWrap(err error, message string) error {
	ae := &applicationError{
		WithStack: NewCustomErrorStack(err, http.StatusInternalServerError, message),
	}

	return ae
}

func NewApplicationErrorWrapWithCode(err error, code int, message string) error {
	ae := &applicationError{
		WithStack: NewCustomErrorStack(err, code, message),
	}

	return ae
}

type applicationError struct {
	contracts.WithStack
}

type ApplicationError interface {
	contracts.WithStack
	GetCustomError() CustomError
	IsApplicationError() bool
}

func (a *applicationError) IsApplicationError() bool {
	return true
}

func (a *applicationError) GetCustomError() CustomError {
	return GetCustomError(a)
}

func IsApplicationError(err error) bool {
	a, ok := err.(ApplicationError)
	if ok && a.IsApplicationError() {
		return true
	}

	var applicationError ApplicationError
	//us, ok := grpc_errors.Cause(err).(ApplicationError)
	if errors.As(err, &applicationError) {
		return applicationError.IsApplicationError()
	}

	return false
}
