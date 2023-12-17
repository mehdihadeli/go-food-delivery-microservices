package customErrors

import (
	"net/http"

	"emperror.dev/errors"
)

func NewInternalServerError(message string) error {
	br := &internalServerError{
		CustomError: NewCustomError(nil, http.StatusInternalServerError, message),
	}
	stackErr := errors.WithStackIf(br)

	return stackErr
}

func NewInternalServerErrorWrap(err error, message string) error {
	br := &internalServerError{
		CustomError: NewCustomError(err, http.StatusInternalServerError, message),
	}
	stackErr := errors.WithStackIf(br)

	return stackErr
}

type internalServerError struct {
	CustomError
}

type InternalServerError interface {
	CustomError
}

func (i *internalServerError) isInternalServerError() bool {
	return true
}

func IsInternalServerError(err error) bool {
	var internalErr *internalServerError
	// us, ok := grpc_errors.Cause(err).(InternalServerError)
	if errors.As(err, &internalErr) {
		return internalErr.isInternalServerError()
	}

	return false
}