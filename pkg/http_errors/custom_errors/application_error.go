package customErrors

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/http_errors/contracts"
	"github.com/pkg/errors"
)

func NewApplicationError(message string, code int) *applicationError {
	ae := &applicationError{
		customError: NewCustomError(nil, code, message),
	}

	return ae
}

func NewApplicationErrorWrap(err error, code int, message string) *applicationError {
	ae := &applicationError{
		customError: NewCustomError(err, code, message),
	}

	return ae
}

type applicationError struct {
	*customError
}

type ApplicationError interface {
	CustomError
	contracts.StackError
	IsApplicationError() bool
}

func (a *applicationError) IsApplicationError() bool {
	return true
}

func (a *applicationError) WithStack() error {
	// with this we use `Cause`, `Unwrap` method of new stack error but this struct `Cause`, `Unwrap` will call with next `Unwrap` on this object
	// Format this error (stackErr) with sprintf, First write Causer of error and then will write call stack for this point of code
	return errors.WithStack(a)
}

func IsApplicationError(err error) bool {
	var applicationError ApplicationError

	//us, ok := errors.Cause(err).(applicationError)
	if errors.As(err, &applicationError) {
		return applicationError.IsApplicationError()
	}

	return false
}
