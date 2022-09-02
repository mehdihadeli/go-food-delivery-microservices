package customErrors

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/http_errors/contracts"
	"github.com/pkg/errors"
)

func NewDomainError(message string, code int) error {
	de := &domainError{
		WithStack: NewCustomErrorStack(nil, code, message),
	}

	return de
}

func NewDomainErrorWrap(err error, code int, message string) error {
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
	d, ok := err.(DomainError)
	if ok && d.IsDomainError() {
		return true
	}

	var domainErr DomainError
	//us, ok := errors.Cause(err).(DomainError)
	if errors.As(err, &domainErr) {
		return domainErr.IsDomainError()
	}

	return false
}
