package dtoV1

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type ProductDto struct {
	ProductId   uuid.UUID `json:"productId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
