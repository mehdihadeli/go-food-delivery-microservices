package updating_product

import (
	uuid "github.com/satori/go.uuid"
)

type UpdateProduct struct {
	ProductID   uuid.UUID `json:"productId" validate:"required,gte=0,lte=255"`
	Name        string    `json:"name" validate:"required,gte=0,lte=255"`
	Description string    `json:"description" validate:"required,gte=0,lte=5000"`
	Price       float64   `json:"price" validate:"required,gte=0"`
}

func NewUpdateProduct(productID uuid.UUID, name string, description string, price float64) UpdateProduct {
	return UpdateProduct{ProductID: productID, Name: name, Description: description, Price: price}
}
