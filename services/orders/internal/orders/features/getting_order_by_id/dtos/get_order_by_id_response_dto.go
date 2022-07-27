package dtos

import "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/dtos"

type GetOrderByIdResponseDto struct {
	Order *dtos.OrderDto `json:"order"`
}
