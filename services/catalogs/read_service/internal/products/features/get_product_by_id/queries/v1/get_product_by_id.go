package v1

import (
	uuid "github.com/satori/go.uuid"
)

type GetProductById struct {
	Id uuid.UUID `validate:"required"`
}
