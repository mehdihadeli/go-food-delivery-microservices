package queries

import (
	uuid "github.com/satori/go.uuid"
)

type GetProductById struct {
	Id uuid.UUID `validate:"required"`
}

// TODO
// this function seems doesn't validate if the id is empty, should I refactor this func and handle the error where NewGetProductByID called?
func NewGetProductById(id uuid.UUID) *GetProductById {
	return &GetProductById{Id: id}
}
