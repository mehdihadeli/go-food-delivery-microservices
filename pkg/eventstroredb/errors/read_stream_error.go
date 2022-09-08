package errors

import (
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	"github.com/pkg/errors"
)

type readStreamError struct {
	customErrors.InternalServerError
}

func NewReadStreamError(err error) error {
	br := &readStreamError{
		InternalServerError: customErrors.NewInternalServerErrorWrap(err, "unable to read events from stream").(customErrors.InternalServerError),
	}

	return br
}

func IsReadStreamError(err error) bool {
	var re *readStreamError
	res := errors.As(err, &re)

	return res
}
