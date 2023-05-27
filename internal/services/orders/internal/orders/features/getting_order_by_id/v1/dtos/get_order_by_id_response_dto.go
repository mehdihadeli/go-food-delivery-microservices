package dtos

import dtosV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orders/internal/orders/dtos/v1"

type GetOrderByIdResponseDto struct {
	Order *dtosV1.OrderReadDto `json:"order"`
}
