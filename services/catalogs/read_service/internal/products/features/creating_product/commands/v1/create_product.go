package v1

import "time"

type CreateProductCommand struct {
	ProductID   string    `validate:"required"`
	Name        string    `validate:"required,min=3,max=250"`
	Description string    `validate:"required,min=3,max=500"`
	Price       float64   `validate:"required"`
	CreatedAt   time.Time `validate:"required"`
}

func NewCreateProductCommand(productID string, name string, description string, price float64, createdAt time.Time) *CreateProductCommand {
	return &CreateProductCommand{ProductID: productID, Name: name, Description: description, Price: price, CreatedAt: createdAt}
}
