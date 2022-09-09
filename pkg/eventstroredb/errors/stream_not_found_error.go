package errors

import (
	"emperror.dev/errors"
	"fmt"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
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
