package v1

import "time"

type CreateProduct struct {
	ProductId   string    `validate:"required"`
	Name        string    `validate:"required,min=3,max=250"`
	Description string    `validate:"required,min=3,max=500"`
	Price       float64   `validate:"required"`
	CreatedAt   time.Time `validate:"required"`
}

func NewCreateProduct(productId string, name string, description string, price float64, createdAt time.Time) *CreateProduct {
	return &CreateProduct{ProductId: productId, Name: name, Description: description, Price: price, CreatedAt: createdAt}
}
