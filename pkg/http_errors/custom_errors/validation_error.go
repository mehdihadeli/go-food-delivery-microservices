package customErrors

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/http_errors/contracts"
	"github.com/pkg/errors"
)

func NewValidationError(message string) error {
	ve := &validationError{
		WithStack: NewBadRequestError(message).(contracts.WithStack),
	}

	return ve
}

func NewValidationErrorWrap(err error, message string) error {
	ve := &validationError{
		WithStack: NewBadRequestErrorWrap(err, message).(contracts.WithStack),
	}

	return ve
}

type validationError struct {
	contracts.WithStack
}

type ValidationError interface {
	BadRequestError
	IsValidationError() bool
}

func (v *validationError) IsValidationError() bool {
	return true
}

func (v *validationError) IsBadRequestError() bool {
	return true
}

func (v *validationError) GetCustomError() CustomError {
	return GetCustomError(v)
}

func IsValidationError(err error) bool {
	v, ok := err.(ValidationError)
	if ok && v.IsValidationError() {
		return true
	}

	var validationError ValidationError
	//us, ok := errors.Cause(err).(ValidationError)
	if errors.As(err, &validationError) {
		return validationError.IsValidationError()
	}

	return false
}
