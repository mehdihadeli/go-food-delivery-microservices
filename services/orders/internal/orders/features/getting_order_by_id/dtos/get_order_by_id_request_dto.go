package dtos

import uuid "github.com/satori/go.uuid"

type GetOrderByIdRequestDto struct {
	OrderId uuid.UUID `param:"id" json:"-"`
}
