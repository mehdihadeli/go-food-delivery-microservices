package customErrors

import (
	"net/http"

	"emperror.dev/errors"
)

func NewMarshalingError(message string) MarshalingError {
	// `NewPlain` doesn't add stack-trace at all
	marshalingErrMessage := errors.NewPlain("marshaling error")
	// `WrapIf` add stack-trace if not added before
	stackErr := errors.WrapIf(marshalingErrMessage, message)

	marshalingError := &marshalingError{
		CustomError: NewCustomError(stackErr, http.StatusInternalServerError, message),
	}

	return marshalingError
}

func NewMarshalingErrorWrap(err error, message string) MarshalingError {
	if err == nil {
		return NewMarshalingError(message)
	}

	// `WithMessage` doesn't add stack-trace at all
	marshalingErrMessage := errors.WithMessage(err, "marshaling error")
	// `WrapIf` add stack-trace if not added before
	stackErr := errors.WrapIf(marshalingErrMessage, message)

	marshalingError := &marshalingError{
		CustomError: NewCustomError(stackErr, http.StatusInternalServerError, message),
	}

	return marshalingError
}

type marshalingError struct {
	CustomError
}

type MarshalingError interface {
	InternalServerError
	isMarshalingError()
}

func (m *marshalingError) isMarshalingError() {
}

func (m *marshalingError) isInternalServerError() {
}

func IsMarshalingError(err error) bool {
	var marshalingErr MarshalingError

	// https://github.com/golang/go/blob/master/src/net/error_windows.go#L10C2-L12C3
	// this doesn't work for a nested marshaling error, and we should use errors.As for traversing errors in all levels
	if _, ok := err.(MarshalingError); ok {
		return true
	}

	// us, ok := errors.Cause(err).(MarshalingError)
	if errors.As(err, &marshalingErr) {
		return true
	}

	return false
}
