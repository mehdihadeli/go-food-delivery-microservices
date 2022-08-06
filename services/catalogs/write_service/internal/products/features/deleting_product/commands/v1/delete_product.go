package v1

import (
	uuid "github.com/satori/go.uuid"
)

type DeleteProductCommand struct {
	ProductID uuid.UUID `validate:"required"`
}

func NewDeleteProductCommand(productID uuid.UUID) *DeleteProductCommand {
	return &DeleteProductCommand{ProductID: productID}
}
