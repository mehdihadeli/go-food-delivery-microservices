package errors

import (
    "fmt"

    "emperror.dev/errors"

    customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
)

type appendToStreamError struct {
    customErrors.BadRequestError
}

type AppendToStreamError interface {
    customErrors.BadRequestError
    IsAppendToStreamError() bool
}

func NewAppendToStreamError(err error, streamId string) error {
    bad := customErrors.NewBadRequestErrorWrap(err, fmt.Sprintf("unable to append events to stream %s", streamId))
    customErr := customErrors.GetCustomError(bad)
    br := &appendToStreamError{
        BadRequestError: customErr.(customErrors.BadRequestError),
    }

    return errors.WithStackIf(br)
}

func (err *appendToStreamError) IsAppendToStreamError() bool {
    return true
}

func IsAppendToStreamError(err error) bool {
    var an AppendToStreamError
    if errors.As(err, &an) {
        return an.IsAppendToStreamError()
    }

    return false
}
