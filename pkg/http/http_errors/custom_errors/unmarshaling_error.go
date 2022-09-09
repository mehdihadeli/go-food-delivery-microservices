package customErrors

import (
	"emperror.dev/errors"
)

func NewUnMarshalingError(message string) error {
	internal := NewInternalServerError(message)
	customErr := GetCustomError(internal)
	ue := &unMarshalingError{
		InternalServerError: customErr.(InternalServerError),
	}
	stackErr := errors.WithStackIf(ue)

	return stackErr
}

func NewUnMarshalingErrorWrap(err error, message string) error {
	internal := NewInternalServerErrorWrap(err, message)
	customErr := GetCustomError(internal)
	ue := &unMarshalingError{
		InternalServerError: customErr.(InternalServerError),
	}
	stackErr := errors.WithStackIf(ue)

	return stackErr
}

type unMarshalingError struct {
	InternalServerError
}

type UnMarshalingError interface {
	InternalServerError
	IsUnMarshalingError() bool
}

func (u *unMarshalingError) IsUnMarshalingError() bool {
	return true
}

func IsUnMarshalingError(err error) bool {
	var unMarshalingError UnMarshalingError
	//us, ok := grpc_errors.Cause(err).(UnMarshalingError)
	if errors.As(err, &unMarshalingError) {
		return unMarshalingError.IsUnMarshalingError()
	}

	return false
}
