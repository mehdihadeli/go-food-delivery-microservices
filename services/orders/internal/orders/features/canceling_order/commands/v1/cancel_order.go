package v1

import (
	uuid "github.com/satori/go.uuid"
)

type CancelOrderCommand struct {
	OrderId      uuid.UUID `json:"orderId" validate:"required"`
	CancelReason string    `json:"cancelReason" validate:"required"`
}

func NewCancelOrderCommand(orderId uuid.UUID, cancelReason string) *CancelOrderCommand {
	return &CancelOrderCommand{OrderId: orderId, CancelReason: cancelReason}
}
