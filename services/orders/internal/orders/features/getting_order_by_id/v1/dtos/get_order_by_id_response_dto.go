package dtos

import "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/dtos/v1"

type GetOrderByIdResponseDto struct {
	Order *dtosV1.OrderReadDto `json:"order"`
}
