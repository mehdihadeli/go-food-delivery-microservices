package commands

import (
	"time"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils/validator"

	uuid "github.com/satori/go.uuid"
)

type UpdateProduct struct {
	ProductID   uuid.UUID `validate:"required"`
	Name        string    `validate:"required,gte=0,lte=255"`
	Description string    `validate:"required,gte=0,lte=5000"`
	Price       float64   `validate:"required,gte=0"`
	UpdatedAt   time.Time `validate:"required"`
}

func NewUpdateProduct(productID uuid.UUID, name string, description string, price float64) (*UpdateProduct, error) {
	command := &UpdateProduct{
		ProductID:   productID,
		Name:        name,
		Description: description,
		Price:       price,
		UpdatedAt:   time.Now(),
	}
	err := validator.Validate(command)
	if err != nil {
		return nil, err
	}

	return command, nil
}
