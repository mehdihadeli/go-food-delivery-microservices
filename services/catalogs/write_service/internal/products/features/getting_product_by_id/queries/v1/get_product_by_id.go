package v1

import (
	uuid "github.com/satori/go.uuid"
)

//https://echo.labstack.com/guide/request/
//https://github.com/go-playground/validator

type GetProductById struct {
	ProductID uuid.UUID `validate:"required"`
}

func NewGetProductById(productId uuid.UUID) *GetProductById {
	return &GetProductById{ProductID: productId}
}
