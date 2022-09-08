package errors

import (
	"fmt"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	"github.com/pkg/errors"
)

type streamNotFoundError struct {
	customErrors.NotFoundError
}

func NewStreamNotFoundError(err error, streamId string) error {
	br := &streamNotFoundError{
		NotFoundError: customErrors.NewNotFoundErrorWrap(err, fmt.Sprintf("stream with streamId %s not found", streamId)).(customErrors.NotFoundError),
	}

	return br
}

func IsStreamNotFoundError(err error) bool {
	var re *streamNotFoundError
	res := errors.As(err, &re)

	return res
}
