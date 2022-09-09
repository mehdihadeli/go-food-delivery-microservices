package customErrors

import (
	"emperror.dev/errors"
	"net/http"
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
	IsUnAuthorizedError() bool
}

func (u *unauthorizedError) IsUnAuthorizedError() bool {
	return true
}

func IsUnAuthorizedError(err error) bool {
	var unauthorizedError UnauthorizedError
	//us, ok := grpc_errors.Cause(err).(UnauthorizedError)
	if errors.As(err, &unauthorizedError) {
		return unauthorizedError.IsUnAuthorizedError()
	}

	return false
}
