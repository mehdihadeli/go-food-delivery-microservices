package domainExceptions

import (
	"fmt"

	customErrors "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/httperrors/customerrors"

	"emperror.dev/errors"
)

type orderNotFoundError struct {
	customErrors.NotFoundError
}

type OrderNotFoundError interface {
	customErrors.NotFoundError
}

func NewOrderNotFoundError(id int) error {
	notFound := customErrors.NewNotFoundError(
		fmt.Sprintf("order with id %d not found", id),
	)
	customErr := customErrors.GetCustomError(notFound).(customErrors.NotFoundError)
	br := &orderNotFoundError{
		NotFoundError: customErr,
	}

	return errors.WithStackIf(br)
}

func (i *orderNotFoundError) isorderNotFoundError() bool {
	return true
}

func IsOrderNotFoundError(err error) bool {
	var os *orderNotFoundError
	if errors.As(err, &os) {
		return os.isorderNotFoundError()
	}

	return false
}
