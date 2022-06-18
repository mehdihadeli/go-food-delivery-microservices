package deleting_product

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/consts"
	uuid "github.com/satori/go.uuid"
)

type DeleteProduct struct {
	ProductID uuid.UUID `json:"productId" validate:"required"`
}

func (DeleteProduct) Key() int { return consts.DeleteProductKey }

func NewDeleteProduct(productID uuid.UUID) *DeleteProduct {
	return &DeleteProduct{ProductID: productID}
}
