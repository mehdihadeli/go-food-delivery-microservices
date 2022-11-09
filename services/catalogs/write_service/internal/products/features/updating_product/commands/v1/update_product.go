package v1

import (
	"time"

	"github.com/go-playground/validator"

	uuid "github.com/satori/go.uuid"
)

type UpdateProduct struct {
	ProductID   uuid.UUID `validate:"required"`
	Name        string    `validate:"required,gte=0,lte=255"`
	Description string    `validate:"required,gte=0,lte=5000"`
	Price       float64   `validate:"required,gte=0"`
	UpdatedAt   time.Time `validate:"required"`
}

func NewUpdateProduct(productID uuid.UUID, name string, description string, price float64) *UpdateProduct {
	return &UpdateProduct{ProductID: productID, Name: name, Description: description, Price: price, UpdatedAt: time.Now()}
}

func (u *UpdateProduct) Validate() error {
	validate := validator.New()

	if err := validate.Struct(u); err != nil {
		// this check ensures there wasn't an error
		// with the validation process itself
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}
		return err
	}
	return nil
}
