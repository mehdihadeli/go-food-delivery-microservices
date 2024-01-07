package customErrors

import (
	"net/http"

	"emperror.dev/errors"
)

func NewBadRequestError(message string) BadRequestError {
	// `NewPlain` doesn't add stack-trace at all
	badRequestErrMessage := errors.NewPlain("bad request error")
	// `WrapIf` add stack-trace if not added before
	stackErr := errors.WrapIf(badRequestErrMessage, message)

	badRequestError := &badRequestError{
		CustomError: NewCustomError(stackErr, http.StatusBadRequest, message),
	}

	return badRequestError
}

func NewBadRequestErrorWrap(err error, message string) BadRequestError {
	if err == nil {
		return NewBadRequestError(message)
	}

	// `WithMessage` doesn't add stack-trace at all
	badRequestErrMessage := errors.WithMessage(err, "bad request error")
	// `WrapIf` add stack-trace if not added before
	stackErr := errors.WrapIf(badRequestErrMessage, message)

	badRequestError := &badRequestError{
		CustomError: NewCustomError(stackErr, http.StatusBadRequest, message),
	}

	return badRequestError
}

type badRequestError struct {
	CustomError
}

type BadRequestError interface {
	CustomError
	isBadRequestError()
}

func (b *badRequestError) isBadRequestError() {
}

func IsBadRequestError(err error) bool {
	var badRequestError BadRequestError

	// https://github.com/golang/go/blob/master/src/net/error_windows.go#L10C2-L12C3
	// this doesn't work for a nested bad-request error, and we should use errors.As for traversing errors in all levels
	if _, ok := err.(BadRequestError); ok {
		return true
	}

	// us, ok := errors.Cause(err).(BadRequestError)
	if errors.As(err, &badRequestError) {
		return true
	}

	return false
}
