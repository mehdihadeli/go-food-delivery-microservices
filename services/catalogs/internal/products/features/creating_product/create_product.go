package creating_product

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/consts"
	uuid "github.com/satori/go.uuid"
)

type CreateProduct struct {
	ProductID   uuid.UUID `json:"productId" validate:"required"`
	Name        string    `json:"name" validate:"required,gte=0,lte=255"`
	Description string    `json:"description" validate:"required,gte=0,lte=5000"`
	Price       float64   `json:"price" validate:"required,gte=0"`
}

func (CreateProduct) Key() int { return consts.CreateProductKey }

func NewCreateProduct(productID uuid.UUID, name string, description string, price float64) *CreateProduct {
	return &CreateProduct{ProductID: productID, Name: name, Description: description, Price: price}
}
