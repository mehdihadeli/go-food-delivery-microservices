package customErrors

import (
	"github.com/pkg/errors"
)

func NewValidationError(message string) error {
	ve := &validationError{
		badRequestError: NewBadRequestError(message),
	}

	return ve
}

func NewValidationErrorWrap(err error, message string) error {
	ve := &validationError{
		badRequestError: NewBadRequestErrorWrap(err, message),
	}

	return ve
}

type validationError struct {
	*badRequestError
}

type ValidationError interface {
	BadRequestError
	IsValidationError() bool
}

func (v *validationError) IsValidationError() bool {
	return true
}

func (v *validationError) WithStack() error {
	// with this we use `Cause`, `Unwrap` method of new stack error but this struct `Cause`, `Unwrap` will call with next `Unwrap` on this object
	// Format this error (stackErr) with sprintf, First write Causer of error and then will write call stack for this point of code
	return errors.WithStack(v)
}

func IsValidationError(err error) bool {
	var validationError ValidationError

	//us, ok := errors.Cause(err).(iBadRequest)
	if errors.As(err, &validationError) {
		return validationError.IsValidationError()
	}

	return false
}
