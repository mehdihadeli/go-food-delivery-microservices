package v1

import (
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/cqrs"
	customErrors "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/httperrors/customerrors"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	uuid "github.com/satori/go.uuid"
)

// https://echo.labstack.com/guide/request/
// https://github.com/go-playground/validator

type GetProductById struct {
	cqrs.Query
	ProductID uuid.UUID
}

func NewGetProductById(productId uuid.UUID) *GetProductById {
	query := &GetProductById{
		Query:     cqrs.NewQueryByT[GetProductById](),
		ProductID: productId,
	}

	return query
}

func NewGetProductByIdWithValidation(productId uuid.UUID) (*GetProductById, error) {
	query := NewGetProductById(productId)
	err := query.Validate()

	return query, err
}

func (p *GetProductById) Validate() error {
	err := validation.ValidateStruct(
		p,
		validation.Field(&p.ProductID, validation.Required, is.UUIDv4),
	)
	if err != nil {
		return customErrors.NewValidationErrorWrap(err, "validation error")
	}

	return nil
}
