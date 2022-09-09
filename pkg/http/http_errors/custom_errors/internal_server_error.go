package customErrors

import (
	"emperror.dev/errors"
	"net/http"
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
	IsInternalServerError() bool
}

func (i *internalServerError) IsInternalServerError() bool {
	return true
}

func IsInternalServerError(err error) bool {
	var internalErr InternalServerError
	//us, ok := grpc_errors.Cause(err).(InternalServerError)
	if errors.As(err, &internalErr) {
		return internalErr.IsInternalServerError()
	}

	return false
}
