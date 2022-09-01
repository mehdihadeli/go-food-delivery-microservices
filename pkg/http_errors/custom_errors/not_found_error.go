package customErrors

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/http_errors/contracts"
	"github.com/pkg/errors"
	"net/http"
)

func NewNotFoundError(message string) *notFoundError {
	ne := &notFoundError{
		customError: NewCustomError(nil, http.StatusNotFound, message),
	}

	return ne
}

func NewNotFoundErrorWrap(err error, message string) *notFoundError {
	ne := &notFoundError{
		customError: NewCustomError(err, http.StatusNotFound, message),
	}

	return ne
}

type notFoundError struct {
	*customError
}

type NotFoundError interface {
	CustomError
	contracts.StackError
	IsNotFoundError() bool
}

func (n *notFoundError) IsNotFoundError() bool {
	return true
}

func (n *notFoundError) WithStack() error {
	return errors.WithStack(n)
}

func IsNotFoundError(err error) bool {
	var notFoundError NotFoundError

	//us, ok := errors.Cause(err).(NotFoundError)
	if errors.As(err, &notFoundError) {
		return notFoundError.IsNotFoundError()
	}

	return false
}
