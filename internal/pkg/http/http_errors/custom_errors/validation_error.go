package customErrors

import (
	"emperror.dev/errors"
)

func NewValidationError(message string) error {
	bad := NewBadRequestError(message)
	customErr := GetCustomError(bad)
	ue := &validationError{
		BadRequestError: customErr.(BadRequestError),
	}
	stackErr := errors.WithStackIf(ue)

	return stackErr
}

func NewValidationErrorWrap(err error, message string) error {
	bad := NewBadRequestErrorWrap(err, message)
	customErr := GetCustomError(bad)
	ue := &validationError{
		BadRequestError: customErr.(BadRequestError),
	}
	stackErr := errors.WithStackIf(ue)

	return stackErr
}

type validationError struct {
	BadRequestError
}

type ValidationError interface {
	BadRequestError
}

func (v *validationError) isValidationError() bool {
	return true
}

func IsValidationError(err error) bool {
	var validationError *validationError
	// us, ok := grpc_errors.Cause(err).(ValidationError)
	if errors.As(err, &validationError) {
		return validationError.isValidationError()
	}

	return false
}
