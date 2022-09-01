package customErrors

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/http_errors/contracts"
	"github.com/pkg/errors"
	"net/http"
)

func NewUnAuthorizedError(message string) *unauthorizedError {
	ue := &unauthorizedError{
		customError: NewCustomError(nil, http.StatusUnauthorized, message),
	}

	return ue
}

func NewUnAuthorizedErrorWrap(err error, message string) *unauthorizedError {
	ue := &unauthorizedError{
		customError: NewCustomError(err, http.StatusUnauthorized, message),
	}

	return ue
}

type unauthorizedError struct {
	*customError
}

type UnauthorizedError interface {
	CustomError
	contracts.StackError
	IsUnAuthorizedError() bool
}

func (u *unauthorizedError) IsUnAuthorizedError() bool {
	return true
}

func (u *unauthorizedError) WithStack() error {
	// with this we use `Cause`, `Unwrap` method of new stack error but this struct `Cause`, `Unwrap` will call with next `Unwrap` on this object
	// Format this error (stackErr) with sprintf, First write Causer of error and then will write call stack for this point of code
	return errors.WithStack(u)
}

func IsUnAuthorizedError(err error) bool {
	var unAuthorizedErr UnauthorizedError

	//us, ok := errors.Cause(err).(UnauthorizedError)
	if errors.As(err, &unAuthorizedErr) {
		return unAuthorizedErr.IsUnAuthorizedError()
	}

	return false
}
