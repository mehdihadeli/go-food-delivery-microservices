package queries

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	uuid "github.com/satori/go.uuid"
)

type GetProductById struct {
	Id uuid.UUID
}

func NewGetProductById(id uuid.UUID) (*GetProductById, error) {
	product := &GetProductById{Id: id}
	if err := product.Validate(); err != nil {
		return nil, err
	}

	return product, nil
}

func (p *GetProductById) Validate() error {
	return validation.ValidateStruct(p, validation.Field(&p.Id, validation.Required, is.UUIDv4))
}
