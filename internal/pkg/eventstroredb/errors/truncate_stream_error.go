package errors

import (
	"fmt"

	customErrors "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/httperrors/customerrors"

	"emperror.dev/errors"
)

type truncateStreamError struct {
	customErrors.InternalServerError
}

type TruncateStreamError interface {
	customErrors.InternalServerError
	IsTruncateStreamError() bool
}

func NewTruncateStreamError(err error, streamId string) error {
	internal := customErrors.NewInternalServerErrorWrap(err, fmt.Sprintf("unable to truncate stream %s", streamId))
	customErr := customErrors.GetCustomError(internal)
	br := &truncateStreamError{
		InternalServerError: customErr.(customErrors.InternalServerError),
	}

	return errors.WithStackIf(br)
}

func (err *truncateStreamError) IsTruncateStreamError() bool {
	return true
}

func IsTruncateStreamError(err error) bool {
	var rs TruncateStreamError
	if errors.As(err, &rs) {
		return rs.IsTruncateStreamError()
	}

	return false
}
