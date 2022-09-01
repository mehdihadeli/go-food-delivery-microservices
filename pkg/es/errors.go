package es

import (
	"fmt"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http_errors/custom_errors"
	"github.com/pkg/errors"
)

var (
	ErrEventAlreadyExists = customErrors.NewConflictError(fmt.Sprintf("domain event already exists in event registry"))
	ErrInvalidEventType   = errors.New("invalid event type")
)
