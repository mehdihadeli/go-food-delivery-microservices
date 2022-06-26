package dtos

import uuid "github.com/satori/go.uuid"

type CreateProductResponseDto struct {
	ProductID uuid.UUID `json:"productId"`
}
