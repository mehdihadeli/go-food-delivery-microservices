package updating_product

import (
	uuid "github.com/satori/go.uuid"
)

type UpdateProduct struct {
	ProductID   uuid.UUID
	Name        string
	Description string
	Price       float64
}

func NewUpdateProduct(productID uuid.UUID, name string, description string, price float64) UpdateProduct {
	return UpdateProduct{ProductID: productID, Name: name, Description: description, Price: price}
}
