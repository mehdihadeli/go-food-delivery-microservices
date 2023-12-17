package v1

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	uuid "github.com/satori/go.uuid"
)

type DeleteProduct struct {
	ProductID uuid.UUID
}

func NewDeleteProduct(productID uuid.UUID) (*DeleteProduct, error) {
	command := &DeleteProduct{ProductID: productID}
	err := command.Validate()
	if err != nil {
		return nil, err
	}

	return command, nil
}

// IsTxRequest for enabling transactions on the mediatr pipeline
func (c *DeleteProduct) IsTxRequest() bool {
	return true
}

func (p *DeleteProduct) Validate() error {
	return validation.ValidateStruct(
		p,
		validation.Field(&p.ProductID, validation.Required),
		validation.Field(&p.ProductID, is.UUIDv4),
	)
}