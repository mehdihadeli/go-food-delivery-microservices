package dtos

import uuid "github.com/satori/go.uuid"

//https://echo.labstack.com/guide/response/

type CreateProductResponseDto struct {
	ProductID uuid.UUID `json:"productId"`
}
