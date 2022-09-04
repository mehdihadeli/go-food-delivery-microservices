package customErrors

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/contracts"
	"github.com/pkg/errors"
	"net/http"
)

func NewForbiddenError(message string) error {
	ne := &forbiddenError{
		WithStack: NewCustomErrorStack(nil, http.StatusForbidden, message),
	}

	return ne
}

func NewForbiddenErrorWrap(err error, message string) error {
	ne := &forbiddenError{
		WithStack: NewCustomErrorStack(err, http.StatusForbidden, message),
	}

	return ne
}

type forbiddenError struct {
	contracts.WithStack
}

type ForbiddenError interface {
	contracts.WithStack
	IsForbiddenError() bool
	GetCustomError() CustomError
}

func (f *forbiddenError) IsForbiddenError() bool {
	return true
}

func (f *forbiddenError) GetCustomError() CustomError {
	return GetCustomError(f)
}

func IsForbiddenError(err error) bool {
	f, ok := err.(ForbiddenError)
	if ok && f.IsForbiddenError() {
		return true
	}

	var forbiddenError ForbiddenError
	//us, ok := grpc_errors.Cause(err).(ForbiddenError)
	if errors.As(err, &forbiddenError) {
		return forbiddenError.IsForbiddenError()
	}

	return false
}
