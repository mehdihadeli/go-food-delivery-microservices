package getProductByIdQuery

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils/validator"
	uuid "github.com/satori/go.uuid"
)

// https://echo.labstack.com/guide/request/
// https://github.com/go-playground/validator

type GetProductById struct {
	ProductID uuid.UUID `validate:"required"`
}

func NewGetProductById(productId uuid.UUID) (*GetProductById, error) {
	query := &GetProductById{ProductID: productId}
	err := validator.Validate(query)
	if err != nil {
		return nil, err
	}

	return query, nil
}
