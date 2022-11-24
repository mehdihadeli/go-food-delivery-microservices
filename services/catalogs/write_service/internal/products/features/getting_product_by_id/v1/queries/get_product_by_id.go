package getProductByIdQuery

import (
	uuid "github.com/satori/go.uuid"

	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils/validator"
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
