package deleting_product

import uuid "github.com/satori/go.uuid"

type DeleteProduct struct {
	ProductID uuid.UUID `json:"productId" validate:"required"`
}

func NewDeleteProductCommand(productID uuid.UUID) *DeleteProduct {
	return &DeleteProduct{ProductID: productID}
}
