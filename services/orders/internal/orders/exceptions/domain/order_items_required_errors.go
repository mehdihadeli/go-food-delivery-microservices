package domain

import (
	"emperror.dev/errors"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
)

type orderShopItemsRequiredError struct {
	customErrors.BadRequestError
}

type OrderShopItemsRequiredError interface {
	customErrors.BadRequestError
	IsOrderShopItemsRequiredError() bool
}

func NewOrderShopItemsRequiredError(message string) error {
	bad := customErrors.NewBadRequestError(message)
	customErr := customErrors.GetCustomError(bad).(customErrors.BadRequestError)
	br := &orderShopItemsRequiredError{
		BadRequestError: customErr,
	}

	return errors.WithStackIf(br)
}

func (err *orderShopItemsRequiredError) IsOrderShopItemsRequiredError() bool {
	return true
}

func IsOrderShopItemsRequiredError(err error) bool {
	var os OrderShopItemsRequiredError
	if errors.As(err, &os) {
		return os.IsOrderShopItemsRequiredError()
	}

	return false
}
