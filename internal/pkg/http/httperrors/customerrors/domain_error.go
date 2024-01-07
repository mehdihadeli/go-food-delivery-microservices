package customErrors

import (
	"net/http"

	"emperror.dev/errors"
)

func NewDomainError(message string) DomainError {
	return NewDomainErrorWithCode(message, http.StatusBadRequest)
}

func NewDomainErrorWithCode(message string, code int) DomainError {
	// `NewPlain` doesn't add stack-trace at all
	domainErrMessage := errors.NewPlain("domain error")
	// `WrapIf` add stack-trace if not added before
	stackErr := errors.WrapIf(domainErrMessage, message)

	domainError := &domainError{
		CustomError: NewCustomError(stackErr, code, message),
	}

	return domainError
}

func NewDomainErrorWrap(err error, message string) DomainError {
	return NewDomainErrorWithCodeWrap(err, http.StatusBadRequest, message)
}

func NewDomainErrorWithCodeWrap(err error, code int, message string) DomainError {
	if err == nil {
		return NewDomainErrorWithCode(message, code)
	}

	// `WithMessage` doesn't add stack-trace at all
	domainErrMessage := errors.WithMessage(err, "domain error")
	// `WrapIf` add stack-trace if not added before
	stackErr := errors.WrapIf(domainErrMessage, message)

	domainError := &domainError{
		CustomError: NewCustomError(stackErr, code, message),
	}

	return domainError
}

type domainError struct {
	CustomError
}

type DomainError interface {
	CustomError
	isDomainError()
}

func (d *domainError) isDomainError() {
}

func IsDomainError(err error, code int) bool {
	var domainErr DomainError

	// https://github.com/golang/go/blob/master/src/net/error_windows.go#L10C2-L12C3
	// this doesn't work for a nested notfound error, and we should use errors.As for traversing errors in all levels
	if _, ok := err.(DomainError); ok {
		return true
	}

	// us, ok := errors.Cause(err).(DomainError)
	if errors.As(err, &domainErr) {
		return domainErr.Status() == code
	}

	return false
}
