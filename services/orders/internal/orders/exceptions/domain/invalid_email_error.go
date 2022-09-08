package domain

import (
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	"github.com/pkg/errors"
)

type InvalidEmailAddressError struct {
	customErrors.BadRequestError
}

func NewInvalidEmailAddressError(message string) error {
	br := &InvalidEmailAddressError{
		BadRequestError: customErrors.NewBadRequestError(message).(customErrors.BadRequestError),
	}

	return br
}

func IsInvalidEmailAddressError(err error) bool {
	var re *InvalidEmailAddressError
	res := errors.As(err, &re)

	return res
}
