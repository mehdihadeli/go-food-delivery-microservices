package creating_product

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type CreateProduct struct {
	ProductID   uuid.UUID `validate:"required"`
	Name        string    `validate:"required,gte=0,lte=255"`
	Description string    `validate:"required,gte=0,lte=5000"`
	Price       float64   `validate:"required,gte=0"`
	CreatedAt   time.Time `validate:"required"`
}

func NewCreateProduct(name string, description string, price float64) *CreateProduct {
	return &CreateProduct{ProductID: uuid.NewV4(), Name: name, Description: description, Price: price, CreatedAt: time.Now()}
}
