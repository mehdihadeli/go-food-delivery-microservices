package dtos

import uuid "github.com/satori/go.uuid"

type GetProductByIdRequestDto struct {
	Id uuid.UUID `param:"id" json:"-"`
}
