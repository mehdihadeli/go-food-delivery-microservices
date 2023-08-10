package commands

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils/validator"
	uuid "github.com/satori/go.uuid"
)

type DeleteProduct struct {
	ProductID uuid.UUID `validate:"required"`
}

func NewDeleteProduct(productID uuid.UUID) (*DeleteProduct, error) {
	command := &DeleteProduct{ProductID: productID}
	err := validator.Validate(command)
	if err != nil {
		return nil, err
	}

	return command, nil
}
