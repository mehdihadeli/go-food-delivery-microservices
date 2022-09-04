package customErrors

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/contracts"
	"github.com/pkg/errors"
)

func NewMarshalingError(message string) error {
	ue := &marshalingError{
		WithStack: NewInternalServerError(message).(contracts.WithStack),
	}

	return ue
}

func NewMarshalingErrorWrap(err error, message string) error {
	ue := &marshalingError{
		WithStack: NewInternalServerErrorWrap(err, message).(contracts.WithStack),
	}

	return ue
}

type marshalingError struct {
	contracts.WithStack
}

type MarshalingError interface {
	InternalServerError
	IsMarshalingError() bool
}

func (m *marshalingError) IsMarshalingError() bool {
	return true
}

func (m *marshalingError) IsInternalServerError() bool {
	return true
}

func (m *marshalingError) GetCustomError() CustomError {
	return GetCustomError(m)
}

func IsMarshalingError(err error) bool {
	m, ok := err.(MarshalingError)
	if ok && m.IsMarshalingError() {
		return true
	}

	var marshalingError MarshalingError
	//us, ok := grpc_errors.Cause(err).(MarshalingError)
	if errors.As(err, &marshalingError) {
		return marshalingError.IsMarshalingError()
	}

	return false
}
