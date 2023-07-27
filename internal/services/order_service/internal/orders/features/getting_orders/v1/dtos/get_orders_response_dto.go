package dtos

import (
	dtosV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/dtos/v1"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils"
)

type GetOrdersResponseDto struct {
	Orders *utils.ListResult[*dtosV1.OrderReadDto]
}
