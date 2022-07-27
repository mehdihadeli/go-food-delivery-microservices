package v1

import (
	uuid "github.com/satori/go.uuid"
)

type SubmitOrderCommand struct {
	OrderId uuid.UUID `validate:"required"`
}

func NewSubmitOrderCommand(orderId uuid.UUID) *SubmitOrderCommand {
	return &SubmitOrderCommand{OrderId: orderId}
}
