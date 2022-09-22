package v1

import uuid "github.com/satori/go.uuid"

type GetOrderById struct {
	Id uuid.UUID `validate:"required"`
}

func NewGetOrderById(id uuid.UUID) *GetOrderById {
	return &GetOrderById{Id: id}
}
