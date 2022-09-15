package dtos

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	ordersDto "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/dtos"
)

type GetOrdersResponseDto struct {
	Orders *utils.ListResult[*ordersDto.OrderReadDto]
}
