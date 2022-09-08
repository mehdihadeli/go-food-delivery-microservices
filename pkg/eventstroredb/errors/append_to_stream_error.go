package errors

import (
	"fmt"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	"github.com/pkg/errors"
)

type appendToStreamError struct {
	customErrors.BadRequestError
}

func NewAppendToStreamError(err error, streamId string) error {
	br := &appendToStreamError{
		BadRequestError: customErrors.NewBadRequestErrorWrap(err, fmt.Sprintf("unable to append events to stream %s", streamId)).(customErrors.BadRequestError),
	}

	return br
}

func IsAppendToStreamError(err error) bool {
	var re *appendToStreamError
	res := errors.As(err, &re)

	return res
}
