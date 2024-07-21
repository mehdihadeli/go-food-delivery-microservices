package domainExceptions

import (
	customErrors "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/httperrors/customerrors"

	"emperror.dev/errors"
)

type invalidDeliveryAddressError struct {
	customErrors.BadRequestError
}
type InvalidDeliveryAddressError interface {
	customErrors.BadRequestError
}

func NewInvalidDeliveryAddressError(message string) error {
	bad := customErrors.NewBadRequestError(message)
	customErr := customErrors.GetCustomError(bad).(customErrors.BadRequestError)
	br := &invalidDeliveryAddressError{
		BadRequestError: customErr,
	}

	return errors.WithStackIf(br)
}

func (i *invalidDeliveryAddressError) isInvalidAddress() bool {
	return true
}

func IsInvalidDeliveryAddressError(err error) bool {
	var ia *invalidDeliveryAddressError
	if errors.As(err, &ia) {
		return ia.isInvalidAddress()
	}

	return false
}
