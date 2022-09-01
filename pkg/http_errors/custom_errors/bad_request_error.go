package customErrors

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/http_errors/contracts"
	"github.com/pkg/errors"
	"net/http"
)

func NewBadRequestError(message string) *badRequestError {
	br := &badRequestError{
		customError: NewCustomError(nil, http.StatusBadRequest, message),
	}

	return br
}

func NewBadRequestErrorWrap(err error, message string) *badRequestError {
	br := &badRequestError{
		customError: NewCustomError(err, http.StatusBadRequest, message),
	}

	return br
}

type badRequestError struct {
	*customError
}

func (b *badRequestError) IsBadRequestError() bool {
	return true
}

func (b *badRequestError) WithStack() error {
	// with this we use `Cause`, `Unwrap` method of new stack error but this struct `Cause`, `Unwrap` will call with next `Unwrap` on this object
	// Format this error (stackErr) with sprintf, First write Causer of error and then will write call stack for this point of code
	return errors.WithStack(b)
}

type BadRequestError interface {
	CustomError
	contracts.StackError
	IsBadRequestError() bool
}

func IsBadRequestError(err error) bool {
	var badErr BadRequestError

	//us, ok := errors.Cause(err).(BadRequestError)
	if errors.As(err, &badErr) {
		return badErr.IsBadRequestError()
	}

	return false
}
