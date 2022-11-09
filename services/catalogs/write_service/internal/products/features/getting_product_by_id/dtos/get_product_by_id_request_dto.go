package dtos

import uuid "github.com/satori/go.uuid"

// https://echo.labstack.com/guide/binding/
// https://echo.labstack.com/guide/request/
// https://github.com/go-playground/validator

// GetProductByIdRequestDto validation will handle in query level
type GetProductByIdRequestDto struct {
	ProductId uuid.UUID `param:"id" json:"-"`
}
