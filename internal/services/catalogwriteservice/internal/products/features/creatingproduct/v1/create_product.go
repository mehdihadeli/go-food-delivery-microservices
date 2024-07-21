package v1

import (
	"time"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/cqrs"
	customErrors "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/httperrors/customerrors"

	validation "github.com/go-ozzo/ozzo-validation"
	uuid "github.com/satori/go.uuid"
)

// https://echo.labstack.com/guide/request/
// https://github.com/go-playground/validator

type CreateProduct struct {
	cqrs.Command
	ProductID   uuid.UUID
	Name        string
	Description string
	Price       float64
	CreatedAt   time.Time
}

// NewCreateProduct Create a new product
func NewCreateProduct(
	name string,
	description string,
	price float64,
) *CreateProduct {
	command := &CreateProduct{
		Command:     cqrs.NewCommandByT[CreateProduct](),
		ProductID:   uuid.NewV4(),
		Name:        name,
		Description: description,
		Price:       price,
		CreatedAt:   time.Now(),
	}

	return command
}

// NewCreateProductWithValidation Create a new product with inline validation - for defensive programming and ensuring validation even without using middleware
func NewCreateProductWithValidation(
	name string,
	description string,
	price float64,
) (*CreateProduct, error) {
	command := NewCreateProduct(name, description, price)
	err := command.Validate()

	return command, err
}

func (c *CreateProduct) isTxRequest(){
}

func (c *CreateProduct) Validate() error {
	err := validation.ValidateStruct(
		c,
		validation.Field(&c.ProductID, validation.Required),
		validation.Field(
			&c.Name,
			validation.Required,
			validation.Length(0, 255),
		),
		validation.Field(
			&c.Description,
			validation.Required,
			validation.Length(0, 5000),
		),
		validation.Field(
			&c.Price,
			validation.Required,
			validation.Min(0.0).Exclusive(),
		),
		validation.Field(&c.CreatedAt, validation.Required),
	)
	if err != nil {
		return customErrors.NewValidationErrorWrap(err, "validation error")
	}

	return nil
}
