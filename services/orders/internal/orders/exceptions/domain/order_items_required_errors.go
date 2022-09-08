package domain

import (
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	"github.com/pkg/errors"
)

type OrderShopItemsRequiredError struct {
	customErrors.BadRequestError
}

func NewOrderShopItemsRequiredError(message string) error {
	br := &OrderShopItemsRequiredError{
		BadRequestError: customErrors.NewBadRequestError(message).(customErrors.BadRequestError),
	}

	return br
}

func IsOrderShopItemsRequiredError(err error) bool {
	var re *OrderShopItemsRequiredError
	res := errors.As(err, &re)

	return res
}
