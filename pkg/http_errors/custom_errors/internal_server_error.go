package customErrors

import (
	"github.com/pkg/errors"
	"net/http"
)

func NewInternalServerError(message string) error {
	br := &internalServerError{
		WithStack: NewCustomErrorStack(nil, http.StatusInternalServerError, message),
	}

	return br
}

func NewInternalServerErrorWrap(err error, message string) error {
	br := &internalServerError{
		WithStack: NewCustomErrorStack(err, http.StatusInternalServerError, message),
	}

	return br
}

type internalServerError struct {
	WithStack
}

func (i *internalServerError) IsInternalServerError() bool {
	return true
}

//func (i *internalServerError) WithStack() error {
//	// with this we use `Cause`, `Unwrap` method of new stack error but this struct `Cause`, `Unwrap` will call with next `Unwrap` on this object
//	// Format this error (stackErr) with sprintf, First write Causer of error and then will write call stack for this point of code
//	return errors.WithStack(i)
//}

type InternalServerError interface {
	WithStack
	IsInternalServerError() bool
}

func IsInternalServerError(err error) bool {
	var internalErr InternalServerError

	_, ok := err.(InternalServerError)
	if ok && internalErr.IsInternalServerError() {
		return true
	}

	//us, ok := errors.Cause(err).(InternalServerError)
	if errors.As(err, &internalErr) {
		return internalErr.IsInternalServerError()
	}

	return false
}
