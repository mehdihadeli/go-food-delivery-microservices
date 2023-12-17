package customErrors

import (
	"net/http"

	"emperror.dev/errors"
)

func NewConflictError(message string) error {
	ce := &conflictError{
		CustomError: NewCustomError(nil, http.StatusConflict, message),
	}
	stackErr := errors.WithStackIf(ce)

	return stackErr
}

func NewConflictErrorWrap(err error, message string) error {
	ce := &conflictError{
		CustomError: NewCustomError(err, http.StatusConflict, message),
	}
	stackErr := errors.WithStackIf(ce)

	return stackErr
}

type conflictError struct {
	CustomError
}

type ConflictError interface {
	CustomError
}

func (c *conflictError) isConflictError() bool {
	return true
}

func IsConflictError(err error) bool {
	var conflictError *conflictError
	// us, ok := grpc_errors.Cause(err).(ConflictError)
	if errors.As(err, &conflictError) {
		return conflictError.isConflictError()
	}

	return false
}
