package v1

import (
	uuid "github.com/satori/go.uuid"
)

type GetProductById struct {
	ProductID uuid.UUID `validate:"required"`
}
