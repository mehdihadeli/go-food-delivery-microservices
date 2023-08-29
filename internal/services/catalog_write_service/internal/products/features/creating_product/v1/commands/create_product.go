package createProductCommand

import (
	validation "github.com/go-ozzo/ozzo-validation"
	uuid "github.com/satori/go.uuid"
	"time"
)

// https://echo.labstack.com/guide/request/
// https://github.com/go-playground/validator

type CreateProduct struct {
	ProductID   uuid.UUID `validate:"required"`
	Name        string    `validate:"required,gte=0,lte=255"`
	Description string    `validate:"required,gte=0,lte=5000"`
	Price       float64   `validate:"required,gte=0"`
	CreatedAt   time.Time `validate:"required"`
}

func NewCreateProduct(name string, description string, price float64) (*CreateProduct, error) {
	command := &CreateProduct{
		ProductID:   uuid.NewV4(),
		Name:        name,
		Description: description,
		Price:       price,
		CreatedAt:   time.Now(),
	}
	err := command.Validate()
	if err != nil {
		return nil, err
	}
	return command, nil
}

func (c *CreateProduct) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.ProductID, validation.Required),
		validation.Field(&c.Name, validation.Required, validation.Length(0, 255)),
		validation.Field(&c.Description, validation.Required, validation.Length(0, 5000)),
		validation.Field(&c.Price, validation.Required, validation.Min(0).Exclusive()),
		validation.Field(&c.CreatedAt, validation.Required),
	)
}
