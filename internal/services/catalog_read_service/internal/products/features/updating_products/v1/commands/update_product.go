package commands

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	uuid "github.com/satori/go.uuid"
)

type UpdateProduct struct {
	ProductId   uuid.UUID
	Name        string
	Description string
	Price       float64
	UpdatedAt   time.Time
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
