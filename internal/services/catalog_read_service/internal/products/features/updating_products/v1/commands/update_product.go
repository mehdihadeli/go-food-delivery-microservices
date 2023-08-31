package commands

import (
	"time"

	"github.com/go-ozzo/ozzo-validation/is"

	validation "github.com/go-ozzo/ozzo-validation"

	uuid "github.com/satori/go.uuid"
)

type UpdateProduct struct {
	ProductId   uuid.UUID `validate:"required"`
	Name        string    `validate:"required,gte=0,lte=255"`
	Description string    `validate:"required,gte=0,lte=5000"`
	Price       float64   `validate:"required,gte=0"`
	UpdatedAt   time.Time `validate:"required"`
}

func NewUpdateProduct(productId uuid.UUID, name string, description string, price float64) (*UpdateProduct, error) {
	product := &UpdateProduct{
		ProductId:   productId,
		Name:        name,
		Description: description,
		Price:       price,
		UpdatedAt:   time.Now(),
	}
	if err := product.Validate(); err != nil {
		return nil, err
	}
	return product, nil
}

func (p *UpdateProduct) Validate() error {
	return validation.ValidateStruct(p, validation.Field(&p.ProductId, validation.Required, is.UUIDv4),
		validation.Field(&p.Name, validation.Required, validation.Length(0, 255)),
		validation.Field(&p.Description, validation.Required, validation.Length(0, 5000)),
		validation.Field(&p.Price, validation.Required, validation.Min(0.0)),
		validation.Field(&p.UpdatedAt, validation.Required),
	)
}
