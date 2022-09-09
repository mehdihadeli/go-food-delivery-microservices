package v1

import uuid "github.com/satori/go.uuid"

type GetOrderByIdQuery struct {
	OrderId uuid.UUID `validate:"required"`
}

func NewGetOrderByIdQuery(id uuid.UUID) *GetOrderByIdQuery {
	return &GetOrderByIdQuery{OrderId: id}
}
