package commands

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	uuid "github.com/satori/go.uuid"
)

type DeleteProduct struct {
	ProductId uuid.UUID
}

func NewDeleteProduct(productId uuid.UUID) (*DeleteProduct, error) {
	delProduct := &DeleteProduct{ProductId: productId}
	if err := delProduct.Validate(); err != nil {
		return nil, err
	}

	return delProduct, nil
}

func (p *DeleteProduct) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.ProductId, validation.Required, is.UUIDv4))
}
