package creating_product

import (
	uuid "github.com/satori/go.uuid"
)

type CreateProduct struct {
	ProductID   uuid.UUID
	Name        string
	Description string
	Price       float64
}

func NewCreateProduct(name string, description string, price float64) CreateProduct {
	return CreateProduct{ProductID: uuid.NewV4(), Name: name, Description: description, Price: price}
}
