package v1

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/dtos"
	uuid "github.com/satori/go.uuid"
)

type PayOrderCommand struct {
	OrderId uuid.UUID       `validate:"required"`
	Payment dtos.PaymentDto `validate:"required"`
}

func NewPayOrderCommand(orderId uuid.UUID, payment dtos.PaymentDto) *PayOrderCommand {
	return &PayOrderCommand{Payment: payment, OrderId: orderId}
}
