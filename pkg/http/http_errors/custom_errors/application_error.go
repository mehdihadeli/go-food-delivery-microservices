package customErrors

import (
	"emperror.dev/errors"
	"net/http"
)

func NewApplicationError(message string) error {
	ae := &applicationError{
		CustomError: NewCustomError(nil, http.StatusInternalServerError, message),
	}
	stackErr := errors.WithStackIf(ae)

	return stackErr
}

func NewApplicationErrorWithCode(message string, code int) error {
	ae := &applicationError{
		CustomError: NewCustomError(nil, code, message),
	}
	stackErr := errors.WithStackIf(ae)

	return stackErr
}

func NewApplicationErrorWrap(err error, message string) error {
	ae := &applicationError{
		CustomError: NewCustomError(err, http.StatusInternalServerError, message),
	}
	stackErr := errors.WithStackIf(ae)

	return stackErr
}

func NewApplicationErrorWrapWithCode(err error, code int, message string) error {
	ae := &applicationError{
		CustomError: NewCustomError(err, code, message),
	}
	stackErr := errors.WithStackIf(ae)

	return stackErr
}

type applicationError struct {
	CustomError
}

type ApplicationError interface {
	CustomError
	IsApplicationError() bool
}

func (a *applicationError) IsApplicationError() bool {
	return true
}

func IsApplicationError(err error) bool {
	var applicationError ApplicationError
	//us, ok := grpc_errors.Cause(err).(ApplicationError)
	if errors.As(err, &applicationError) {
		return applicationError.IsApplicationError()
	}

	return false
}
