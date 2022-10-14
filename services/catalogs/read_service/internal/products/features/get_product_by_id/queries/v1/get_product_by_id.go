package v1

import (
	uuid "github.com/satori/go.uuid"
)

type GetProductById struct {
	Id uuid.UUID `validate:"required"`
}

func NewGetProductById(id uuid.UUID) *GetProductById {
	return &GetProductById{Id: id}
}
