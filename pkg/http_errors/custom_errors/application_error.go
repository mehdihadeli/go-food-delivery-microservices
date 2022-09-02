package customErrors

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/http_errors/contracts"
	"github.com/pkg/errors"
)

func NewApplicationError(message string, code int) error {
	ae := &applicationError{
		WithStack: NewCustomErrorStack(nil, code, message),
	}

	return ae
}

func NewApplicationErrorWrap(err error, code int, message string) error {
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
	//us, ok := errors.Cause(err).(ApplicationError)
	if errors.As(err, &applicationError) {
		return applicationError.IsApplicationError()
	}

	return false
}
