package commands

import (
	uuid "github.com/satori/go.uuid"
)

type SubmitOrder struct {
	OrderId uuid.UUID `validate:"required"`
}

func NewSubmitOrder(orderId uuid.UUID) *SubmitOrder {
	return &SubmitOrder{OrderId: orderId}
}
