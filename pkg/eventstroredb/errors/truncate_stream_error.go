package errors

import (
	"fmt"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	"github.com/pkg/errors"
)

type truncateStreamError struct {
	customErrors.InternalServerError
}

func NewTruncateStreamError(err error, streamId string) error {
	br := &truncateStreamError{
		InternalServerError: customErrors.NewInternalServerErrorWrap(err, fmt.Sprintf("unable to truncate stream %s", streamId)).(customErrors.InternalServerError),
	}

	return br
}

func IsTruncateStreamError(err error) bool {
	var re *truncateStreamError
	res := errors.As(err, &re)

	return res
}
