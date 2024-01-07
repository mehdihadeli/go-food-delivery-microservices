package customErrors

import (
	"emperror.dev/errors"
)

func NewApiError(message string, code int) ApiError {
	// `NewPlain` doesn't add stack-trace at all
	apiErrMessage := errors.NewPlain("api error")
	// `WrapIf` add stack-trace if not added before
	stackErr := errors.WrapIf(apiErrMessage, message)

	apiError := &apiError{
		CustomError: NewCustomError(stackErr, code, message),
	}

	return apiError
}

func NewApiErrorWrap(err error, code int, message string) ApiError {
	if err == nil {
		return NewApiError(message, code)
	}

	// `WithMessage` doesn't add stack-trace at all
	apiErrMessage := errors.WithMessage(err, "api error")
	// `WrapIf` add stack-trace if not added before
	stackErr := errors.WrapIf(apiErrMessage, message)

	apiError := &apiError{
		CustomError: NewCustomError(stackErr, code, message),
	}

	return apiError
}

type apiError struct {
	CustomError
}

type ApiError interface {
	CustomError
	isAPIError()
}

func (a *apiError) isAPIError() {
}

func IsApiError(err error, code int) bool {
	var apiError ApiError

	// https://github.com/golang/go/blob/master/src/net/error_windows.go#L10C2-L12C3
	// this doesn't work for a nested api error, and we should use errors.As for traversing errors in all levels
	if _, ok := err.(ApiError); ok {
		return true
	}

	// us, ok := errors.Cause(err).(ApiError)
	if errors.As(err, &apiError) {
		return apiError.Status() == code
	}

	return false
}
