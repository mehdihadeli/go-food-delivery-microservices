package customErrors

import (
	"net/http"

	"emperror.dev/errors"
)

func NewUnAuthorizedError(message string) error {
	ue := &unauthorizedError{
		CustomError: NewCustomError(nil, http.StatusUnauthorized, message),
	}
	stackErr := errors.WithStackIf(ue)

	return stackErr
}

func NewUnAuthorizedErrorWrap(err error, message string) error {
	ue := &unauthorizedError{
		CustomError: NewCustomError(err, http.StatusUnauthorized, message),
	}
	stackErr := errors.WithStackIf(ue)

	return stackErr
}

type unauthorizedError struct {
	CustomError
}

type UnauthorizedError interface {
	CustomError
}

func (u *unauthorizedError) isUnAuthorizedError() bool {
	return true
}

func IsUnAuthorizedError(err error) bool {
	var unauthorizedError *unauthorizedError
	// us, ok := grpc_errors.Cause(err).(UnauthorizedError)
	if errors.As(err, &unauthorizedError) {
		return unauthorizedError.isUnAuthorizedError()
	}

	return false
}
