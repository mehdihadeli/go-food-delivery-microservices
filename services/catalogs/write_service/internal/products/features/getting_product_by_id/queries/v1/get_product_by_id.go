package v1

import (
	uuid "github.com/satori/go.uuid"
)

type GetProductByIdQuery struct {
	ProductID uuid.UUID `validate:"required"`
}
