package customErrors

import (
	"net/http"

	"emperror.dev/errors"
)

func NewNotFoundError(message string) NotFoundError {
	// `NewPlain` doesn't add stack-trace at all
	notFoundErrMessage := errors.NewPlain("not found error")
	// `WrapIf` add stack-trace if not added before
	stackErr := errors.WrapIf(notFoundErrMessage, message)

	notFoundError := &notFoundError{
		CustomError: NewCustomError(stackErr, http.StatusBadRequest, message),
	}

	return notFoundError
}

func NewNotFoundErrorWrap(err error, message string) NotFoundError {
	if err == nil {
		return NewNotFoundError(message)
	}

	// `WithMessage` doesn't add stack-trace at all
	notFoundErrMessage := errors.WithMessage(err, "not found error")
	// `WrapIf` add stack-trace if not added before
	stackErr := errors.WrapIf(notFoundErrMessage, message)

	notFoundError := &notFoundError{
		CustomError: NewCustomError(stackErr, http.StatusNotFound, message),
	}

	return notFoundError
}

type notFoundError struct {
	CustomError
}

type NotFoundError interface {
	CustomError
	isNotFoundError()
}

func (n *notFoundError) isNotFoundError() {
}

func IsNotFoundError(err error) bool {
	var notFoundError NotFoundError

	// https://github.com/golang/go/blob/master/src/net/error_windows.go#L10C2-L12C3
	// this doesn't work for a nested notfound error, and we should use errors.As for traversing errors in all levels
	if _, ok := err.(NotFoundError); ok {
		return true
	}

	// us, ok := errors.Cause(err).(NotFoundError)
	if errors.As(err, &notFoundError) {
		return true
	}

	return false
}
