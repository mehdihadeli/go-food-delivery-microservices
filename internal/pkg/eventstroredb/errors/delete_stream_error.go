package errors

import (
	"fmt"

	customErrors "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/httperrors/customerrors"

	"emperror.dev/errors"
)

type deleteStreamError struct {
	customErrors.InternalServerError
}

type DeleteStreamError interface {
	customErrors.InternalServerError
	IsDeleteStreamError() bool
}

func NewDeleteStreamError(err error, streamId string) error {
	internal := customErrors.NewInternalServerErrorWrap(err, fmt.Sprintf("unable to delete stream %s", streamId))
	customErr := customErrors.GetCustomError(internal)

	br := &deleteStreamError{
		InternalServerError: customErr.(customErrors.InternalServerError),
	}

	return errors.WithStackIf(br)
}

func (err *deleteStreamError) IsDeleteStreamError() bool {
	return true
}

func IsDeleteStreamError(err error) bool {
	var ds DeleteStreamError
	if errors.As(err, &ds) {
		return ds.IsDeleteStreamError()
	}

	return false
}
