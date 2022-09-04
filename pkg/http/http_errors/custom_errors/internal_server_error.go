package customErrors

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/contracts"
	"github.com/pkg/errors"
	"net/http"
)

func NewInternalServerError(message string) error {
	br := &internalServerError{
		WithStack: NewCustomErrorStack(nil, http.StatusInternalServerError, message),
	}

	return br
}

func NewInternalServerErrorWrap(err error, message string) error {
	br := &internalServerError{
		WithStack: NewCustomErrorStack(err, http.StatusInternalServerError, message),
	}

	return br
}

type internalServerError struct {
	contracts.WithStack
}

func (i *internalServerError) IsInternalServerError() bool {
	return true
}

func (i *internalServerError) GetCustomError() CustomError {
	return GetCustomError(i)
}

type InternalServerError interface {
	contracts.WithStack
	IsInternalServerError() bool
	GetCustomError() CustomError
}

func IsInternalServerError(err error) bool {
	i, ok := err.(InternalServerError)
	if ok && i.IsInternalServerError() {
		return true
	}

	var internalErr InternalServerError
	//us, ok := grpc_errors.Cause(err).(InternalServerError)
	if errors.As(err, &internalErr) {
		return internalErr.IsInternalServerError()
	}

	return false
}
