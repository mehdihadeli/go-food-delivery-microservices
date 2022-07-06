package creating_product

import "time"

type CreateProduct struct {
	ProductID   string  `validate:"required"`
	Name        string  `validate:"required,min=3,max=250"`
	Description string  `validate:"required,min=3,max=500"`
	Price       float64 `validate:"required"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewCreateProduct(productID string, name string, description string, price float64, createdAt time.Time, updatedAt time.Time) CreateProduct {
	return CreateProduct{ProductID: productID, Name: name, Description: description, Price: price, CreatedAt: createdAt, UpdatedAt: updatedAt}
}
