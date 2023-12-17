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
}

func (a *apiError) isAPIError() bool {
	return true
}

func IsApiError(err error, code int) bool {
	var apiError *apiError

	// us, ok := grpc_errors.Cause(err).(ApiError)
	if errors.As(err, &apiError) {
		return apiError.isAPIError() && apiError.Status() == code
	}

	return false
}
