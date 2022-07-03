package updating_product

import uuid "github.com/satori/go.uuid"

// https://echo.labstack.com/guide/binding/

type UpdateProductRequestDto struct {
	ProductID   uuid.UUID `param:"id"  json:"-" validate:"required"`
	Name        string    `json:"name" validate:"required,gte=0,lte=255"`
	Description string    `json:"description" validate:"required,gte=0,lte=5000"`
	Price       float64   `json:"price" validate:"required,gte=0"`
}
