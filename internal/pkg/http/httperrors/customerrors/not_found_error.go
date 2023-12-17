package customErrors

import (
	"net/http"

	"emperror.dev/errors"
)

func NewNotFoundError(message string) error {
	ne := &notFoundError{
		CustomError: NewCustomError(nil, http.StatusNotFound, message),
	}
	stackErr := errors.WithStackIf(ne)

	return stackErr
}

func NewNotFoundErrorWrap(err error, message string) error {
	ne := &notFoundError{
		CustomError: NewCustomError(err, http.StatusNotFound, message),
	}
	stackErr := errors.WithStackIf(ne)

	return stackErr
}

type notFoundError struct {
	CustomError
}

type NotFoundError interface {
	CustomError
}

func (n *notFoundError) isNotFoundError() bool {
	return true
}

func IsNotFoundError(err error) bool {
	var notFoundError *notFoundError
	// us, ok := grpc_errors.Cause(err).(NotFoundError)
	if errors.As(err, &notFoundError) {
		return notFoundError.isNotFoundError()
	}

	return false
}
