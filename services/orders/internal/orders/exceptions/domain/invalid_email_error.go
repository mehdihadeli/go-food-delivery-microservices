package domain

import (
	"emperror.dev/errors"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
)

type invalidEmailAddressError struct {
	customErrors.BadRequestError
}

type InvalidEmailAddressError interface {
	customErrors.BadRequestError
	IsInvalidEmailAddressError() bool
}

func NewInvalidEmailAddressError(message string) error {
	bad := customErrors.NewBadRequestError(message)
	customErr := customErrors.GetCustomError(bad).(customErrors.BadRequestError)
	br := &invalidEmailAddressError{
		BadRequestError: customErr,
	}

	return errors.WithStackIf(br)
}

func (err *invalidEmailAddressError) IsInvalidEmailAddressError() bool {
	return true
}

func IsInvalidEmailAddressError(err error) bool {
	var ie InvalidEmailAddressError
	if errors.As(err, &ie) {
		return ie.IsInvalidEmailAddressError()
	}

	return false
}
