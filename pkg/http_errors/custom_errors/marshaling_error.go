package customErrors

import "github.com/pkg/errors"

func NewMarshalingError(message string) *marshalingError {
	ue := &marshalingError{
		internalServerError: NewInternalServerError(message),
	}

	return ue
}

func NewMarshalingErrorWrap(err error, message string) *marshalingError {
	ue := &marshalingError{
		internalServerError: NewInternalServerErrorWrap(err, message),
	}

	return ue
}

type marshalingError struct {
	internalServerError error
}

type MarshalingError interface {
	InternalServerError
	IsMarshalingError() bool
}

func (m *marshalingError) IsMarshalingError() bool {
	return true
}

//func (m *marshalingError) WithStack() error {
//	// with this we use `Cause`, `Unwrap` method of new stack error but this struct `Cause`, `Unwrap` will call with next `Unwrap` on this object
//	// Format this error (stackErr) with sprintf, First write Causer of error and then will write call stack for this point of code
//	return errors.WithStack(m)
//}

func IsMarshalingError(err error) bool {
	var marshalingError MarshalingError

	//us, ok := errors.Cause(err).(MarshalingError)
	if errors.As(err, &marshalingError) {
		return marshalingError.IsMarshalingError()
	}

	return false
}
