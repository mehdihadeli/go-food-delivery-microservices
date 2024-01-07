package customErrors

import (
	"net/http"

	"emperror.dev/errors"
)

func NewApplicationError(message string) ApplicationError {
	return NewApplicationErrorWithCode(message, http.StatusInternalServerError)
}

func NewApplicationErrorWithCode(message string, code int) ApplicationError {
	// `NewPlain` doesn't add stack-trace at all
	applicationErrMessage := errors.NewPlain("application error")
	// `WrapIf` add stack-trace if not added before
	stackErr := errors.WrapIf(applicationErrMessage, message)

	applicationError := &applicationError{
		CustomError: NewCustomError(stackErr, code, message),
	}

	return applicationError
}

func NewApplicationErrorWrap(err error, message string) ApplicationError {
	return NewApplicationErrorWrapWithCode(err, http.StatusInternalServerError, message)
}

func NewApplicationErrorWrapWithCode(
	err error,
	code int,
	message string,
) ApplicationError {
	if err == nil {
		return NewApplicationErrorWithCode(message, code)
	}

	// `WithMessage` doesn't add stack-trace at all
	applicationErrMessage := errors.WithMessage(err, "application error")
	// `WrapIf` add stack-trace if not added before
	stackErr := errors.WrapIf(applicationErrMessage, message)

	applicationError := &applicationError{
		CustomError: NewCustomError(stackErr, code, message),
	}

	return applicationError
}

type applicationError struct {
	CustomError
}

type ApplicationError interface {
	CustomError
	isApplicationError()
}

func (a *applicationError) isApplicationError() {
}

func IsApplicationError(err error, code int) bool {
	var applicationError ApplicationError

	// https://github.com/golang/go/blob/master/src/net/error_windows.go#L10C2-L12C3
	// this doesn't work for a nested application error, and we should use errors.As for traversing errors in all levels
	if _, ok := err.(ApplicationError); ok {
		return true
	}

	// us, ok := errors.Cause(err).(ApplicationError)
	if errors.As(err, &applicationError) {
		return applicationError.Status() == code
	}

	return false
}
