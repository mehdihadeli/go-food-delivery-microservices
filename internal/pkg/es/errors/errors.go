package errors

import (
	"fmt"

	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"

	"emperror.dev/errors"
)

var (
	EventAlreadyExistsError = customErrors.NewConflictError(
		fmt.Sprintf("domain_events event already exists in event registry"),
	)
	InvalidEventTypeError = errors.New("invalid event type")
)
