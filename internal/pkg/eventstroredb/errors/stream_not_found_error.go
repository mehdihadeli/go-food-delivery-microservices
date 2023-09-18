package errors

import (
	"fmt"

	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"

	"emperror.dev/errors"
)

type streamNotFoundError struct {
	customErrors.NotFoundError
}

type StreamNotFoundError interface {
	customErrors.NotFoundError
	IsStreamNotFoundError() bool
}

func NewStreamNotFoundError(err error, streamId string) error {
	notFound := customErrors.NewNotFoundErrorWrap(err, fmt.Sprintf("stream with streamId %s not found", streamId))
	customErr := customErrors.GetCustomError(notFound)
	br := &streamNotFoundError{
		NotFoundError: customErr.(customErrors.NotFoundError),
	}

	return errors.WithStackIf(br)
}

func (err *streamNotFoundError) IsStreamNotFoundError() bool {
	return true
}

func IsStreamNotFoundError(err error) bool {
	var rs StreamNotFoundError
	if errors.As(err, &rs) {
		return rs.IsStreamNotFoundError()
	}

	return false
}
