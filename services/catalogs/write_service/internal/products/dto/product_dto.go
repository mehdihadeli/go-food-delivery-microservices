package dto

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type ProductDto struct {
	ProductID   uuid.UUID `json:"productId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
