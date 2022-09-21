package v1

import uuid "github.com/satori/go.uuid"

type DeleteProduct struct {
	ProductId uuid.UUID `validate:"required"`
}

func NewDeleteProduct(productId uuid.UUID) *DeleteProduct {
	return &DeleteProduct{ProductId: productId}
}
