package errors

import (
	"fmt"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	"github.com/pkg/errors"
)

type deleteStreamError struct {
	customErrors.InternalServerError
}

func NewDeleteStreamError(err error, streamId string) error {
	br := &deleteStreamError{
		InternalServerError: customErrors.NewInternalServerErrorWrap(err, fmt.Sprintf("unable to delete stream %s", streamId)).(customErrors.InternalServerError),
	}

	return br
}

func IsDeleteStreamError(err error) bool {
	var re *deleteStreamError
	res := errors.As(err, &re)

	return res
}
