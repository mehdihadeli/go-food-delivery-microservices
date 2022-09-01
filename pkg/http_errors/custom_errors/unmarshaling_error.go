package customErrors

import (
	"github.com/pkg/errors"
)

func NewUnMarshalingError(message string) *unMarshalingError {
	ue := &unMarshalingError{
		internalServerError: NewInternalServerError(message),
	}

	return ue
}

func NewUnMarshalingErrorWrap(err error, message string) *unMarshalingError {
	ue := &unMarshalingError{
		internalServerError: NewInternalServerErrorWrap(err, message),
	}

	return ue
}

type unMarshalingError struct {
	internalServerError error
}

type UnMarshalingError interface {
	InternalServerError
	IsUnMarshalingError() bool
}

func (u *unMarshalingError) IsUnMarshalingError() bool {
	return true
}

//func (u *unMarshalingError) WithStack() error {
//	// with this we use `Cause`, `Unwrap` method of new stack error but this struct `Cause`, `Unwrap` will call with next `Unwrap` on this object
//	// Format this error (stackErr) with sprintf, First write Causer of error and then will write call stack for this point of code
//	return errors.WithStack(u)
//}

func IsUnMarshalingError(err error) bool {
	var unMarshalingError UnMarshalingError

	//us, ok := errors.Cause(err).(UnMarshalingError)
	if errors.As(err, &unMarshalingError) {
		return unMarshalingError.IsUnMarshalingError()
	}

	return false
}
