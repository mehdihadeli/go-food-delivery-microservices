package customErrors

import (
	"net/http"

	"emperror.dev/errors"
)

func NewConflictError(message string) ConflictError {
	// `NewPlain` doesn't add stack-trace at all
	conflictErrMessage := errors.NewPlain("conflict error")
	// `WrapIf` add stack-trace if not added before
	stackErr := errors.WrapIf(conflictErrMessage, message)

	conflictError := &conflictError{
		CustomError: NewCustomError(stackErr, http.StatusConflict, message),
	}

	return conflictError
}

func NewConflictErrorWrap(err error, message string) ConflictError {
	if err == nil {
		return NewConflictError(message)
	}

	// `WithMessage` doesn't add stack-trace at all
	conflictErrMessage := errors.WithMessage(err, "conflict error")
	// `WrapIf` add stack-trace if not added before
	stackErr := errors.WrapIf(conflictErrMessage, message)

	conflictError := &conflictError{
		CustomError: NewCustomError(stackErr, http.StatusConflict, message),
	}

	return conflictError
}

type conflictError struct {
	CustomError
}

type ConflictError interface {
	CustomError
	isConflictError()
}

func (c *conflictError) isConflictError() {
}

func IsConflictError(err error) bool {
	var conflictError ConflictError

	// https://github.com/golang/go/blob/master/src/net/error_windows.go#L10C2-L12C3
	// this doesn't work for a nested notfound error, and we should use errors.As for traversing errors in all levels
	if _, ok := err.(ConflictError); ok {
		return true
	}

	// us, ok := errors.Cause(err).(ConflictError)
	if errors.As(err, &conflictError) {
		return true
	}

	return false
}
