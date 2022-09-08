package domain

import (
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	"github.com/pkg/errors"
)

type InvalidDeliveryAddressError struct {
	customErrors.BadRequestError
}

func NewInvalidDeliveryAddressError(message string) error {
	br := &InvalidDeliveryAddressError{
		BadRequestError: customErrors.NewBadRequestError(message).(customErrors.BadRequestError),
	}

	return br
}

func IsInvalidDeliveryAddressError(err error) bool {
	var re *InvalidDeliveryAddressError
	res := errors.As(err, &re)

	return res
}
