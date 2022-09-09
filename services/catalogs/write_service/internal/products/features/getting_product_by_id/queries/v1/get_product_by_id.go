package v1

import (
	uuid "github.com/satori/go.uuid"
)

//https://echo.labstack.com/guide/request/
//https://github.com/go-playground/validator

type GetProductByIdQuery struct {
	ProductID uuid.UUID `validate:"required"`
}

func NewGetProductByIdQuery(productId uuid.UUID) *GetProductByIdQuery {
	return &GetProductByIdQuery{ProductID: productId}
}
