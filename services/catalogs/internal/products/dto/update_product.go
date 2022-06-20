package dto

import uuid "github.com/satori/go.uuid"

type UpdateProductRequestDto struct {
	ProductID   uuid.UUID `json:"productId" validate:"required,gte=0,lte=255"`
	Name        string    `json:"name" validate:"required,gte=0,lte=255"`
	Description string    `json:"description" validate:"required,gte=0,lte=5000"`
	Price       float64   `json:"price" validate:"required,gte=0"`
}

type UpdateProductResponseDto struct {
	ProductID   uuid.UUID `json:"productId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
}
