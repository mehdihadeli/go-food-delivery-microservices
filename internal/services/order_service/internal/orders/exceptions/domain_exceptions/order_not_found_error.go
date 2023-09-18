package domainExceptions

import (
	"fmt"

	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"

	"emperror.dev/errors"
)

type orderNotFoundError struct {
	customErrors.NotFoundError
}

type OrderNotFoundError interface {
	customErrors.NotFoundError
}

func NewOrderNotFoundError(id int) error {
	notFound := customErrors.NewNotFoundError(fmt.Sprintf("order with id %d not found", id))
	customErr := customErrors.GetCustomError(notFound).(customErrors.NotFoundError)
	br := &orderNotFoundError{
		NotFoundError: customErr,
	}

	return errors.WithStackIf(br)
}

func IsOrderNotFoundError(err error) bool {
	var os OrderNotFoundError
	if errors.As(err, &os) {
		return true
	}

	return false
}
