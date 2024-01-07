package customErrors

import (
	"net/http"

	"emperror.dev/errors"
)

func NewForbiddenError(message string) ForbiddenError {
	// `NewPlain` doesn't add stack-trace at all
	forbiddenErrMessage := errors.NewPlain("forbidden error")
	// `WrapIf` add stack-trace if not added before
	stackErr := errors.WrapIf(forbiddenErrMessage, message)

	forbiddenError := &forbiddenError{
		CustomError: NewCustomError(stackErr, http.StatusForbidden, message),
	}

	return forbiddenError
}

func NewForbiddenErrorWrap(err error, message string) ForbiddenError {
	if err == nil {
		return NewForbiddenError(message)
	}

	// `WithMessage` doesn't add stack-trace at all
	forbiddenErrMessage := errors.WithMessage(err, "forbidden error")
	// `WrapIf` add stack-trace if not added before
	stackErr := errors.WrapIf(forbiddenErrMessage, message)

	forbiddenError := &forbiddenError{
		CustomError: NewCustomError(stackErr, http.StatusForbidden, message),
	}

	return forbiddenError
}

type forbiddenError struct {
	CustomError
}

type ForbiddenError interface {
	CustomError
	isForbiddenError()
}

func (f *forbiddenError) isForbiddenError() {
}

func IsForbiddenError(err error) bool {
	var forbiddenError ForbiddenError

	// https://github.com/golang/go/blob/master/src/net/error_windows.go#L10C2-L12C3
	// this doesn't work for a nested forbidden error, and we should use errors.As for traversing errors in all levels
	if _, ok := err.(ForbiddenError); ok {
		return true
	}

	// us, ok := errors.Cause(err).(ForbiddenError)
	if errors.As(err, &forbiddenError) {
		return true
	}

	return false
}
