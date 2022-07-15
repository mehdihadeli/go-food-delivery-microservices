package dtos

import uuid "github.com/satori/go.uuid"

type GetProductByIdRequestDto struct {
	ProductId uuid.UUID `param:"id" json:"-"`
}
