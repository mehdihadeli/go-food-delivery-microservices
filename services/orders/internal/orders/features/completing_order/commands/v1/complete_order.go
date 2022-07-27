package v1

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type CompleteOrderCommand struct {
	OrderId           uuid.UUID `validate:"required"`
	DeliveryTimestamp time.Time `validate:"required"`
}

func NewCompleteOrderCommand(orderId uuid.UUID, deliveryTimestamp time.Time) *CompleteOrderCommand {
	return &CompleteOrderCommand{OrderId: orderId, DeliveryTimestamp: deliveryTimestamp}
}
