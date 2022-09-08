package domain

import (
	"fmt"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	"github.com/pkg/errors"
)

type OrderNotFoundError struct {
	customErrors.NotFoundError
}

func NewOrderNotFoundError(id int) error {
	br := &OrderNotFoundError{
		NotFoundError: customErrors.NewNotFoundError(fmt.Sprintf("order with id %d not found", id)).(customErrors.NotFoundError),
	}

	return br
}

func IsOrderNotFoundError(err error) bool {
	var re *OrderNotFoundError
	res := errors.As(err, &re)

	return res
}
