package commands

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	uuid "github.com/satori/go.uuid"
)

type UpdateProduct struct {
	ProductID   uuid.UUID
	Name        string
	Description string
	Price       float64
	UpdatedAt   time.Time
}

func NewUpdateProduct(productID uuid.UUID, name string, description string, price float64) (*UpdateProduct, error) {
	command := &UpdateProduct{
		ProductID:   productID,
		Name:        name,
		Description: description,
		Price:       price,
		UpdatedAt:   time.Now(),
	}
	err := command.Validate()
	if err != nil {
		return nil, err
	}

	return command, nil
}

func (p *UpdateProduct) Validate() error {
	return validation.ValidateStruct(p, validation.Field(&p.ProductID, validation.Required),
		validation.Field(&p.Name, validation.Required, validation.Length(0, 255)),
		validation.Field(&p.Description, validation.Required, validation.Length(0, 5000)),
		validation.Field(&p.Price, validation.Required, validation.Min(0.0)),
		validation.Field(&p.UpdatedAt, validation.Required))
}
