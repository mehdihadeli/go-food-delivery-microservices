package integration_events

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
	dtoV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/dto/v1"
	uuid "github.com/satori/go.uuid"
)

type ProductCreatedV1 struct {
	*types.Message
	*dtoV1.ProductDto
}

func NewProductCreatedV1(productDto *dtoV1.ProductDto) *ProductCreatedV1 {
	return &ProductCreatedV1{ProductDto: productDto, Message: types.NewMessage(uuid.NewV4().String())}
}
