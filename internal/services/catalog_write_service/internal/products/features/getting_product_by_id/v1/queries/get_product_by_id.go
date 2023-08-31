package getProductByIdQuery

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	uuid "github.com/satori/go.uuid"
)

// https://echo.labstack.com/guide/request/
// https://github.com/go-playground/validator

type GetProductById struct {
	ProductID uuid.UUID `validate:"required"`
}

func NewGetProductById(productId uuid.UUID) (*GetProductById, error) {
	query := &GetProductById{ProductID: productId}
	err := query.Validate()
	if err != nil {
		return nil, err
	}

	return query, nil
}

func (p *GetProductById) Validate() error {
	return validation.ValidateStruct(p, validation.Field(&p.ProductID, validation.Required, is.UUIDv4))
}
