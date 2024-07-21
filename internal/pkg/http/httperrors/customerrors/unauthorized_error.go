package customErrors

import (
	"net/http"

	"emperror.dev/errors"
)

func NewUnAuthorizedError(message string) UnauthorizedError {
	// `NewPlain` doesn't add stack-trace at all
	unAuthorizedErrMessage := errors.NewPlain("unauthorized error")
	// `WrapIf` add stack-trace if not added before
	stackErr := errors.WrapIf(unAuthorizedErrMessage, message)

	unAuthorizedError := &unauthorizedError{
		CustomError: NewCustomError(stackErr, http.StatusUnauthorized, message),
	}

	return unAuthorizedError
}

func NewUnAuthorizedErrorWrap(err error, message string) UnauthorizedError {
	if err == nil {
		return NewUnAuthorizedError(message)
	}

	// `WithMessage` doesn't add stack-trace at all
	unAuthorizedErrMessage := errors.WithMessage(err, "unauthorized error")
	// `WrapIf` add stack-trace if not added before
	stackErr := errors.WrapIf(unAuthorizedErrMessage, message)

	unAuthorizedError := &unauthorizedError{
		CustomError: NewCustomError(stackErr, http.StatusUnauthorized, message),
	}

	return unAuthorizedError
}

type unauthorizedError struct {
	CustomError
}

type UnauthorizedError interface {
	CustomError
	isUnAuthorizedError()
}

func (u *unauthorizedError) isUnAuthorizedError() {
}

func IsUnAuthorizedError(err error) bool {
	var unauthorizedError UnauthorizedError

	// https://github.com/golang/go/blob/master/src/net/error_windows.go#L10C2-L12C3
	// this doesn't work for a nested unauthorized error, and we should use errors.As for traversing errors in all levels
	if _, ok := err.(UnauthorizedError); ok {
		return true
	}

	// us, ok := errors.Cause(err).(UnauthorizedError)
	if errors.As(err, &unauthorizedError) {
		return true
	}

	return false
}
