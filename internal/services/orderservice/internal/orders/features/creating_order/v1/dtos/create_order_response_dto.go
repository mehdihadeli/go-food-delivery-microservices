package dtos

import uuid "github.com/satori/go.uuid"

// https://echo.labstack.com/guide/response/
type CreateOrderResponseDto struct {
	OrderId uuid.UUID `json:"Id"`
}
