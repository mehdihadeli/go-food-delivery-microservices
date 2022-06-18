package getting_product_by_id

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/consts"
	uuid "github.com/satori/go.uuid"
)

type GetProductById struct {
	ProductID uuid.UUID `json:"productId" validate:"required,gte=0,lte=255"`
}

func (GetProductById) Key() int { return consts.GetProductByIdKey }

func NewGetProductById(productID uuid.UUID) *GetProductById {
	return &GetProductById{ProductID: productID}
}
