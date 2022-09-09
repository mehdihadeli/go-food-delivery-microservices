package customErrors

import (
	"emperror.dev/errors"
	"net/http"
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
	IsConflictError() bool
}

func (c *conflictError) IsConflictError() bool {
	return true
}

func IsConflictError(err error) bool {
	var conflictError ConflictError
	//us, ok := grpc_errors.Cause(err).(ConflictError)
	if errors.As(err, &conflictError) {
		return conflictError.IsConflictError()
	}

	return false
}
