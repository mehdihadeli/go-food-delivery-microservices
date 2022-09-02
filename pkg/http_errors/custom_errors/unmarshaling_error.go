package customErrors

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/http_errors/contracts"
	"github.com/pkg/errors"
)

func NewUnMarshalingError(message string) error {
	ue := &marshalingError{
		WithStack: NewInternalServerError(message).(contracts.WithStack),
	}

	return ue
}

func NewUnMarshalingErrorWrap(err error, message string) error {
	ue := &marshalingError{
		WithStack: NewInternalServerErrorWrap(err, message).(contracts.WithStack),
	}

	return ue
}

type unMarshalingError struct {
	contracts.WithStack
}

type UnMarshalingError interface {
	InternalServerError
	IsUnMarshalingError() bool
}

func (u *unMarshalingError) IsUnMarshalingError() bool {
	return true
}

func (u *unMarshalingError) IsInternalServerError() bool {
	return true
}

func (u *unMarshalingError) GetCustomError() CustomError {
	return GetCustomError(u)
}

func IsUnMarshalingError(err error) bool {
	u, ok := err.(UnMarshalingError)
	if ok && u.IsUnMarshalingError() {
		return true
	}

	var unMarshalingError UnMarshalingError
	//us, ok := errors.Cause(err).(UnMarshalingError)
	if errors.As(err, &unMarshalingError) {
		return unMarshalingError.IsUnMarshalingError()
	}

	return false
}
