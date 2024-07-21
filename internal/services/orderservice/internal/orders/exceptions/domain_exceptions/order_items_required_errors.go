package domainExceptions

import (
	customErrors "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/httperrors/customerrors"

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

func (i *orderShopItemsRequiredError) isOrderShopItemsRequiredError() bool {
	return true
}

func IsOrderShopItemsRequiredError(err error) bool {
	var os *orderShopItemsRequiredError
	if errors.As(err, &os) {
		return os.isOrderShopItemsRequiredError()
	}

	return false
}
