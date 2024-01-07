package customErrors

import (
	"net/http"

	"emperror.dev/errors"
)

func NewUnMarshalingError(message string) UnMarshalingError {
	// `NewPlain` doesn't add stack-trace at all
	unMarshalingErrMessage := errors.NewPlain("unMarshaling error")
	// `WrapIf` add stack-trace if not added before
	stackErr := errors.WrapIf(unMarshalingErrMessage, message)

	unMarshalingError := &unMarshalingError{
		CustomError: NewCustomError(stackErr, http.StatusInternalServerError, message),
	}

	return unMarshalingError
}

func NewUnMarshalingErrorWrap(err error, message string) UnMarshalingError {
	if err == nil {
		return NewUnMarshalingError(message)
	}

	// `WithMessage` doesn't add stack-trace at all
	unMarshalingErrMessage := errors.WithMessage(err, "unMarshaling error")
	// `WrapIf` add stack-trace if not added before
	stackErr := errors.WrapIf(unMarshalingErrMessage, message)

	unMarshalingError := &unMarshalingError{
		CustomError: NewCustomError(stackErr, http.StatusInternalServerError, message),
	}

	return unMarshalingError
}

type unMarshalingError struct {
	CustomError
}

type UnMarshalingError interface {
	InternalServerError
	isUnMarshalingError()
}

func (u *unMarshalingError) isUnMarshalingError() {
}

func (u *unMarshalingError) isInternalServerError() {
}

func IsUnMarshalingError(err error) bool {
	var unMarshalingError UnMarshalingError

	// https://github.com/golang/go/blob/master/src/net/error_windows.go#L10C2-L12C3
	// this doesn't work for a nested unMarshaling error, and we should use errors.As for traversing errors in all levels
	if _, ok := err.(UnMarshalingError); ok {
		return true
	}

	// us, ok := errors.Cause(err).(UnMarshalingError)
	if errors.As(err, &unMarshalingError) {
		return true
	}

	return false
}
