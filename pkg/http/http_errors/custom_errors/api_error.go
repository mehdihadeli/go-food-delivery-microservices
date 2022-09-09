package customErrors

import (
	"emperror.dev/errors"
)

func NewApiError(message string, code int) error {
	ae := &apiError{
		CustomError: NewCustomError(nil, code, message),
	}
	stackErr := errors.WithStackIf(ae)

	return stackErr
}

func NewApiErrorWrap(err error, code int, message string) error {
	ae := &apiError{
		CustomError: NewCustomError(err, code, message),
	}
	stackErr := errors.WithStackIf(ae)

	return stackErr
}

type apiError struct {
	CustomError
}

type ApiError interface {
	CustomError
	IsApiError() bool
}

func (a *apiError) IsApiError() bool {
	return true
}

func IsApiError(err error) bool {
	var apiError ApiError
	//us, ok := grpc_errors.Cause(err).(ApiError)
	if errors.As(err, &apiError) {
		return apiError.IsApiError()
	}

	return false
}
