package customErrors

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/contracts"
	"github.com/pkg/errors"
	"net/http"
)

func NewBadRequestError(message string) error {
	br := &badRequestError{
		WithStack: NewCustomErrorStack(nil, http.StatusBadRequest, message),
	}

	return br
}

func NewBadRequestErrorWrap(err error, message string) error {
	br := &badRequestError{
		WithStack: NewCustomErrorStack(err, http.StatusBadRequest, message),
	}

	return br
}

type badRequestError struct {
	contracts.WithStack
}

type BadRequestError interface {
	contracts.WithStack
	IsBadRequestError() bool
	GetCustomError() CustomError
}

func (b *badRequestError) IsBadRequestError() bool {
	return true
}

func (b *badRequestError) GetCustomError() CustomError {
	return GetCustomError(b)
}

func IsBadRequestError(err error) bool {
	var badRequestError *badRequestError
	//us, ok := grpc_errors.Cause(err).(*badRequestError)
	if errors.As(err, &badRequestError) {
		return badRequestError.IsBadRequestError()
	}

	return false
}
