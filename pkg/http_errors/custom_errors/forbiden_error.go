package customErrors

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/http_errors/contracts"
	"github.com/pkg/errors"
	"net/http"
)

func NewNotForbiddenError(message string) *forbiddenError {
	ne := &forbiddenError{
		customError: NewCustomError(nil, http.StatusNotFound, message),
	}

	// Every time an error is wrapped, callers() is called. This seems wasteful to me. I think the root error's stack trace is the most important and almost always what I want. (https://github.com/pkg/errors/issues/75)
	// with this we use `Cause`, `Unwrap` method of new stack error but this struct `Cause`, `Unwrap` will call with next `Unwrap` on this object
	return ne
}

func NewNotForbiddenErrorWrap(err error, message string) *forbiddenError {
	ne := &forbiddenError{
		customError: NewCustomError(err, http.StatusNotFound, message),
	}

	return ne
}

type forbiddenError struct {
	*customError
}

func (f *forbiddenError) IsForbiddenError() bool {
	return true
}

func (f *forbiddenError) WithStack() error {
	// with this we use `Cause`, `Unwrap` method of new stack error but this struct `Cause`, `Unwrap` will call with next `Unwrap` on this object
	// Format this error (stackErr) with sprintf, First write Causer of error and then will write call stack for this point of code
	return errors.WithStack(f)
}

type ForbiddenError interface {
	CustomError
	contracts.StackError
	IsForbiddenError() bool
}

func IsForbiddenError(err error) bool {
	var forbiddenErr ForbiddenError

	//us, ok := errors.Cause(err).(ForbiddenError)
	if errors.As(err, &forbiddenErr) {
		return forbiddenErr.IsForbiddenError()
	}

	return false
}
