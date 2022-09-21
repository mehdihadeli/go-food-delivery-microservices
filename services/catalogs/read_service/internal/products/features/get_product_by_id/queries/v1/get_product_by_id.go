package v1

import (
	uuid "github.com/satori/go.uuid"
)

type GetProductById struct {
	ProductId uuid.UUID `validate:"required"`
}
