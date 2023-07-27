package integration_events

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/types"
	uuid "github.com/satori/go.uuid"

	dtoV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/dto/v1"
)

type ProductCreatedV1 struct {
	*types.Message
	*dtoV1.ProductDto
}

func NewProductCreatedV1(productDto *dtoV1.ProductDto) *ProductCreatedV1 {
	return &ProductCreatedV1{ProductDto: productDto, Message: types.NewMessage(uuid.NewV4().String())}
}
