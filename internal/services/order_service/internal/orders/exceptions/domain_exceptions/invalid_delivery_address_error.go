package domainExceptions

import (
	"emperror.dev/errors"
	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
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

func IsInvalidDeliveryAddressError(err error) bool {
	var ia InvalidDeliveryAddressError
	if errors.As(err, &ia) {
		return true
	}

	return false
}
