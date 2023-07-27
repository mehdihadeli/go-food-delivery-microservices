package domainExceptions

import (
    "fmt"

    "emperror.dev/errors"
    customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
)

type orderNotFoundError struct {
	customErrors.NotFoundError
}

type OrderNotFoundError interface {
	customErrors.NotFoundError
	IsOrderNotFoundError() bool
}

func NewOrderNotFoundError(id int) error {
	notFound := customErrors.NewNotFoundError(fmt.Sprintf("order with id %d not found", id))
	customErr := customErrors.GetCustomError(notFound).(customErrors.NotFoundError)
	br := &orderNotFoundError{
		NotFoundError: customErr,
	}

	return errors.WithStackIf(br)
}

func (err *orderNotFoundError) IsOrderNotFoundError() bool {
	return true
}

func IsOrderNotFoundError(err error) bool {
	var os OrderNotFoundError
	if errors.As(err, &os) {
		return os.IsOrderNotFoundError()
	}

	return false
}
