package getting_product_by_id

import (
	uuid "github.com/satori/go.uuid"
)

type GetProductById struct {
	ProductID uuid.UUID `json:"productId" validate:"required,gte=0,lte=255"`
}

func NewGetProductById(productID uuid.UUID) *GetProductById {
	return &GetProductById{ProductID: productID}
}
