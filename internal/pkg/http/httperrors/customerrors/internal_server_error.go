package customErrors

import (
	"net/http"

	"emperror.dev/errors"
)

func NewInternalServerError(message string) InternalServerError {
	// `NewPlain` doesn't add stack-trace at all
	internalErrMessage := errors.NewPlain("internal server error")
	// `WrapIf` add stack-trace if not added before
	stackErr := errors.WrapIf(internalErrMessage, message)

	internalServerError := &internalServerError{
		CustomError: NewCustomError(stackErr, http.StatusInternalServerError, message),
	}

	return internalServerError
}

func NewInternalServerErrorWrap(err error, message string) InternalServerError {
	if err == nil {
		return NewInternalServerError(message)
	}

	// `WithMessage` doesn't add stack-trace at all
	internalErrMessage := errors.WithMessage(err, "internal server error")
	// `WrapIf` add stack-trace if not added before
	stackErr := errors.WrapIf(internalErrMessage, message)

	internalServerError := &internalServerError{
		CustomError: NewCustomError(stackErr, http.StatusInternalServerError, message),
	}

	return internalServerError
}

type internalServerError struct {
	CustomError
}

type InternalServerError interface {
	CustomError
	isInternalServerError()
}

func (i *internalServerError) isInternalServerError() {
}

func IsInternalServerError(err error) bool {
	var internalServerErr InternalServerError

	// https://github.com/golang/go/blob/master/src/net/error_windows.go#L10C2-L12C3
	// this doesn't work for a nested internal server error, and we should use errors.As for traversing errors in all levels
	if _, ok := err.(InternalServerError); ok {
		return true
	}

	// us, ok := errors.Cause(err).(InternalServerError)
	if errors.As(err, &internalServerErr) {
		return true
	}

	return false
}
