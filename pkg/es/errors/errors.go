package errors

import (
	"fmt"

	"emperror.dev/errors"

	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
)

var (
	EventAlreadyExistsError = customErrors.NewConflictError(fmt.Sprintf("domain_events event already exists in event registry"))
	InvalidEventTypeError   = errors.New("invalid event type")
)
