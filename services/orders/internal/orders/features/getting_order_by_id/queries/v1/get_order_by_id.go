package v1

import uuid "github.com/satori/go.uuid"

type GetOrderById struct {
	OrderId uuid.UUID `validate:"required"`
}

func NewGetOrderById(id uuid.UUID) *GetOrderById {
	return &GetOrderById{OrderId: id}
}
