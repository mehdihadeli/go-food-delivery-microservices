package customErrors

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/contracts"
	"github.com/pkg/errors"
	"net/http"
)

func NewConflictError(message string) error {
	ce := &conflictError{
		WithStack: NewCustomErrorStack(nil, http.StatusConflict, message),
	}

	return ce
}

func NewConflictErrorWrap(err error, message string) error {
	ce := &conflictError{
		WithStack: NewCustomErrorStack(err, http.StatusConflict, message),
	}

	return ce
}

type conflictError struct {
	contracts.WithStack
}

type ConflictError interface {
	contracts.WithStack
	GetCustomError() CustomError
	IsConflictError() bool
}

func (c *conflictError) IsConflictError() bool {
	return true
}

func (c *conflictError) GetCustomError() CustomError {
	return GetCustomError(c)
}

func IsConflictError(err error) bool {
	c, ok := err.(ConflictError)
	if ok && c.IsConflictError() {
		return true
	}

	var conflictError ConflictError
	//us, ok := grpc_errors.Cause(err).(ConflictError)
	if errors.As(err, &conflictError) {
		return conflictError.IsConflictError()
	}

	return false
}
