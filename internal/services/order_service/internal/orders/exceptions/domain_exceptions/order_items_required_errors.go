package domainExceptions

import (
	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"

	"emperror.dev/errors"
)

type orderShopItemsRequiredError struct {
	customErrors.BadRequestError
}

type OrderShopItemsRequiredError interface {
	customErrors.BadRequestError
}

func NewOrderShopItemsRequiredError(message string) error {
	bad := customErrors.NewBadRequestError(message)
	customErr := customErrors.GetCustomError(bad).(customErrors.BadRequestError)
	br := &orderShopItemsRequiredError{
		BadRequestError: customErr,
	}

	return errors.WithStackIf(br)
}

func IsOrderShopItemsRequiredError(err error) bool {
	var os OrderShopItemsRequiredError
	if errors.As(err, &os) {
		return true
	}

	return false
}
