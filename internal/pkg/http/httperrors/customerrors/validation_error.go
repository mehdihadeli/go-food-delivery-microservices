package customErrors

import (
	"net/http"

	"emperror.dev/errors"
)

func NewValidationError(message string) ValidationError {
	// `NewPlain` doesn't add stack-trace at all
	validationErrMessage := errors.NewPlain("validation error")
	// `WrapIf` add stack-trace if not added before
	stackErr := errors.WrapIf(validationErrMessage, message)

	validationError := &validationError{
		CustomError: NewCustomError(stackErr, http.StatusBadRequest, message),
	}

	return validationError
}

func NewValidationErrorWrap(err error, message string) ValidationError {
	if err == nil {
		return NewValidationError(message)
	}

	// `WithMessage` doesn't add stack-trace at all
	validationErrMessage := errors.WithMessage(err, "validation error")
	// `WrapIf` add stack-trace if not added before
	stackErr := errors.WrapIf(validationErrMessage, message)

	validationError := &validationError{
		CustomError: NewCustomError(stackErr, http.StatusBadRequest, message),
	}

	return validationError
}

type validationError struct {
	CustomError
}

type ValidationError interface {
	BadRequestError
	isValidationError()
}

func (v *validationError) isValidationError() {
}

func (v *validationError) isBadRequestError() {
}

func IsValidationError(err error) bool {
	var validationError ValidationError

	// https://github.com/golang/go/blob/master/src/net/error_windows.go#L10C2-L12C3
	// this doesn't work for a nested validation error, and we should use errors.As for traversing errors in all levels
	if _, ok := err.(ValidationError); ok {
		return true
	}

	// us, ok := errors.Cause(err).(ValidationError)
	if errors.As(err, &validationError) {
		return true
	}

	return false
}
