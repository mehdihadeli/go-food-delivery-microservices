package customErrors

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/http_errors/contracts"
	"github.com/pkg/errors"
)

func NewDomainError(message string, code int) *domainError {
	de := &domainError{
		customError: NewCustomError(nil, code, message),
	}

	return de
}

func NewDomainErrorWrap(err error, code int, message string) *domainError {
	de := &domainError{
		customError: NewCustomError(err, code, message),
	}

	return de
}

type domainError struct {
	*customError
}

type DomainError interface {
	CustomError
	contracts.StackError
	IsDomainError() bool
}

func (d *domainError) IsDomainError() bool {
	return true
}

func (d *domainError) WithStack() error {
	// with this we use `Cause`, `Unwrap` method of new stack error but this struct `Cause`, `Unwrap` will call with next `Unwrap` on this object
	// Format this error (stackErr) with sprintf, First write Causer of error and then will write call stack for this point of code
	return errors.WithStack(d)
}

func IsDomainError(err error) bool {
	var domainErr DomainError

	//us, ok := errors.Cause(err).(DomainError)
	if errors.As(err, &domainErr) {
		return domainErr.IsDomainError()
	}

	return false
}
