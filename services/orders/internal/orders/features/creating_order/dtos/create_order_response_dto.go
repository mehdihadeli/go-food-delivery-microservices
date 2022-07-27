package dtos

import uuid "github.com/satori/go.uuid"

type CreateOrderResponseDto struct {
	OrderID uuid.UUID `json:"orderID"`
}
