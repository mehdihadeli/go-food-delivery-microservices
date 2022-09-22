package dtos

import uuid "github.com/satori/go.uuid"

type GetOrderByIdRequestDto struct {
	Id uuid.UUID `param:"id" json:"-"`
}
