package errors

import (
	"fmt"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

type aggregateNotFoundError struct {
	customErrors.NotFoundError
}

func NewAggregateNotFoundError(err error, id uuid.UUID) error {
	br := &aggregateNotFoundError{
		NotFoundError: customErrors.NewNotFoundErrorWrap(err, fmt.Sprintf("aggregtae with id %s not found", id.String())).(customErrors.NotFoundError),
	}

	return br
}

func IsAggregateNotFoundError(err error) bool {
	var re *aggregateNotFoundError
	res := errors.As(err, &re)

	return res
}
