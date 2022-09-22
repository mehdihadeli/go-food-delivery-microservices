package v1

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/dtos"
	uuid "github.com/satori/go.uuid"
)

type OrderCreatedV1 struct {
	*types.Message
	*dtos.OrderReadDto
}

func NewOrderCreatedV1(orderReadDto *dtos.OrderReadDto) *OrderCreatedV1 {
	return &OrderCreatedV1{OrderReadDto: orderReadDto, Message: types.NewMessage(uuid.NewV4().String())}
}
