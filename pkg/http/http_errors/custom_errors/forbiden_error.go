package customErrors

import (
	"emperror.dev/errors"
	"net/http"
)

func NewForbiddenError(message string) error {
	ne := &forbiddenError{
		CustomError: NewCustomError(nil, http.StatusForbidden, message),
	}
	stackErr := errors.WithStackIf(ne)

	return stackErr
}

func NewForbiddenErrorWrap(err error, message string) error {
	ne := &forbiddenError{
		CustomError: NewCustomError(err, http.StatusForbidden, message),
	}
	stackErr := errors.WithStackIf(ne)

	return stackErr
}

type forbiddenError struct {
	CustomError
}

type ForbiddenError interface {
	CustomError
	IsForbiddenError() bool
}

func (f *forbiddenError) IsForbiddenError() bool {
	return true
}

func IsForbiddenError(err error) bool {
	var forbiddenError ForbiddenError
	//us, ok := grpc_errors.Cause(err).(ForbiddenError)
	if errors.As(err, &forbiddenError) {
		return forbiddenError.IsForbiddenError()
	}

	return false
}
