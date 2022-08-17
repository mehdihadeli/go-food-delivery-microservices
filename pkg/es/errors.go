package es

import (
	"fmt"
	httpErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http_errors"
	"github.com/pkg/errors"
)

var (
	ErrEventAlreadyExists = httpErrors.NewConflictError(nil, fmt.Sprintf("domain event already exists in event registry"))
	ErrInvalidEventType   = errors.New("invalid event type")
)
