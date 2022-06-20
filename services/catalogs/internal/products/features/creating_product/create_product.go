package creating_product

import (
	uuid "github.com/satori/go.uuid"
)

type CreateProduct struct {
	ProductID   uuid.UUID `json:"productId" validate:"required"`
	Name        string    `json:"name" validate:"required,gte=0,lte=255"`
	Description string    `json:"description" validate:"required,gte=0,lte=5000"`
	Price       float64   `json:"price" validate:"required,gte=0"`
}

func NewCreateProduct(name string, description string, price float64) *CreateProduct {
	return &CreateProduct{ProductID: uuid.NewV4(), Name: name, Description: description, Price: price}
}
