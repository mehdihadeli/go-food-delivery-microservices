package domainExceptions

import (
	"emperror.dev/errors"
	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
)

type invalidEmailAddressError struct {
	customErrors.BadRequestError
}

type InvalidEmailAddressError interface {
	customErrors.BadRequestError
}

func NewInvalidEmailAddressError(message string) error {
	bad := customErrors.NewBadRequestError(message)
	customErr := customErrors.GetCustomError(bad).(customErrors.BadRequestError)
	br := &invalidEmailAddressError{
		BadRequestError: customErr,
	}

	return errors.WithStackIf(br)
}

func IsInvalidEmailAddressError(err error) bool {
	var ie InvalidEmailAddressError
	if errors.As(err, &ie) {
		return true
	}

	return false
}
