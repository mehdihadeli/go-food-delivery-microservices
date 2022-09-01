package customErrors

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/http_errors/contracts"
	"github.com/pkg/errors"
	"net/http"
)

func NewConflictError(message string) *conflictError {
	ce := &conflictError{
		customError: NewCustomError(nil, http.StatusConflict, message),
	}

	return ce
}

func NewConflictErrorWrap(err error, message string) *conflictError {
	ce := &conflictError{
		customError: NewCustomError(err, http.StatusConflict, message),
	}

	return ce
}

type conflictError struct {
	*customError
}

func (c *conflictError) IsConflictError() bool {
	return true
}

func (c *conflictError) WithStack() error {
	// with this we use `Cause`, `Unwrap` method of new stack error but this struct `Cause`, `Unwrap` will call with next `Unwrap` on this object
	// Format this error (stackErr) with sprintf, First write Causer of error and then will write call stack for this point of code
	return errors.WithStack(c)
}

type ConflictError interface {
	CustomError
	contracts.StackError
	IsConflictError() bool
}

func IsConflictError(err error) bool {
	var conflictError ConflictError

	//us, ok := errors.Cause(err).(ConflictError)
	if errors.As(err, &conflictError) {
		return conflictError.IsConflictError()
	}

	return false
}
