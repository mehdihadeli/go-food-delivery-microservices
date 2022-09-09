package domain

import (
	"emperror.dev/errors"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
)

type invalidDeliveryAddressError struct {
	customErrors.BadRequestError
}
type InvalidDeliveryAddressError interface {
	customErrors.BadRequestError
	IsInvalidDeliveryAddressError() bool
}

func NewInvalidDeliveryAddressError(message string) error {
	bad := customErrors.NewBadRequestError(message)
	customErr := customErrors.GetCustomError(bad).(customErrors.BadRequestError)
	br := &invalidDeliveryAddressError{
		BadRequestError: customErr,
	}

	return errors.WithStackIf(br)
}

func (err *invalidDeliveryAddressError) IsInvalidDeliveryAddressError() bool {
	return true
}

func IsInvalidDeliveryAddressError(err error) bool {
	var ia InvalidDeliveryAddressError
	if errors.As(err, &ia) {
		return ia.IsInvalidDeliveryAddressError()
	}

	return false
}
