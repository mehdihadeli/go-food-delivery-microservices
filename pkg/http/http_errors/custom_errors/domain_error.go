package customErrors

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/contracts"
	"github.com/pkg/errors"
	"net/http"
)

func NewDomainError(message string) error {
	de := &domainError{
		WithStack: NewCustomErrorStack(nil, http.StatusBadRequest, message),
	}

	return de
}

func NewDomainErrorWithCode(message string, code int) error {
	de := &domainError{
		WithStack: NewCustomErrorStack(nil, code, message),
	}

	return de
}

func NewDomainErrorWrap(err error, message string) error {
	de := &domainError{
		WithStack: NewCustomErrorStack(err, http.StatusBadRequest, message),
	}

	return de
}

func NewDomainErrorWithCodeWrap(err error, code int, message string) error {
	de := &domainError{
		WithStack: NewCustomErrorStack(err, code, message),
	}

	return de
}

type domainError struct {
	contracts.WithStack
}

type DomainError interface {
	contracts.WithStack
	GetCustomError() CustomError
	IsDomainError() bool
}

func (d *domainError) IsDomainError() bool {
	return true
}

func (d *domainError) GetCustomError() CustomError {
	return GetCustomError(d)
}

func IsDomainError(err error) bool {
	var domainErr *domainError
	//us, ok := grpc_errors.Cause(err).(*domainError)
	if errors.As(err, &domainErr) {
		return domainErr.IsDomainError()
	}

	return false
}
