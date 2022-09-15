package errors

import (
	"emperror.dev/errors"
	"fmt"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
)

var (
	EventAlreadyExistsError = customErrors.NewConflictError(fmt.Sprintf("domain event already exists in event registry"))
	InvalidEventTypeError   = errors.New("invalid event type")
)
