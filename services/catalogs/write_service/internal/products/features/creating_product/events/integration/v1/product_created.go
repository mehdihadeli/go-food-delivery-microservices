package v1

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/dto"
	uuid "github.com/satori/go.uuid"
)

type ProductCreatedV1 struct {
	*types.Message
	*dto.ProductDto
}

func NewProductCreatedV1(productDto *dto.ProductDto) *ProductCreatedV1 {
	return &ProductCreatedV1{ProductDto: productDto, Message: types.NewMessage(uuid.NewV4().String())}
}
