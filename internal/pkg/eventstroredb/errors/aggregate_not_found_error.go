package errors

import (
    "fmt"

    "emperror.dev/errors"
    uuid "github.com/satori/go.uuid"

    customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
)

//https://klotzandrew.com/blog/error-handling-in-golang/
//https://banzaicloud.com/blog/error-handling-go/

type aggregateNotFoundError struct {
	customErrors.NotFoundError
}

type AggregateNotFoundError interface {
	customErrors.NotFoundError
	IsAggregateNotFoundError() bool
}

func NewAggregateNotFoundError(err error, id uuid.UUID) error {
	notFound := customErrors.NewNotFoundErrorWrap(err, fmt.Sprintf("aggregtae with id %s not found", id.String()))
	customErr := customErrors.GetCustomError(notFound)
	br := &aggregateNotFoundError{
		NotFoundError: customErr.(customErrors.NotFoundError),
	}

	return errors.WithStackIf(br)
}

func (err *aggregateNotFoundError) IsAggregateNotFoundError() bool {
	return true
}

func IsAggregateNotFoundError(err error) bool {
	var an AggregateNotFoundError
	if errors.As(err, &an) {
		return an.IsAggregateNotFoundError()
	}

	return false
}
