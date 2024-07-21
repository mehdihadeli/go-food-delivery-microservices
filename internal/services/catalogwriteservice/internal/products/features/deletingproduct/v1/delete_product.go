package v1

import (
	customErrors "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/httperrors/customerrors"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	uuid "github.com/satori/go.uuid"
)

type DeleteProduct struct {
	ProductID uuid.UUID
}

// NewDeleteProduct delete a product
func NewDeleteProduct(productID uuid.UUID) *DeleteProduct {
	command := &DeleteProduct{ProductID: productID}

	return command
}

// NewDeleteProductWithValidation delete a product with inline validation - for defensive programming and ensuring validation even without using middleware
func NewDeleteProductWithValidation(productID uuid.UUID) (*DeleteProduct, error) {
	command := NewDeleteProduct(productID)
	err := command.Validate()

	return command, err
}

// IsTxRequest for enabling transactions on the mediatr pipeline
func (c *DeleteProduct) isTxRequest() {
}

func (c *DeleteProduct) Validate() error {
	err := validation.ValidateStruct(
		c,
		validation.Field(&c.ProductID, validation.Required),
		validation.Field(&c.ProductID, is.UUIDv4),
	)
	if err != nil {
		return customErrors.NewValidationErrorWrap(err, "validation error")
	}

	return nil
}
