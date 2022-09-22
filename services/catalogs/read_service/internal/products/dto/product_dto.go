package dto

import (
	"time"
)

type ProductDto struct {
	Id          string    `json:"id"`
	ProductId   string    `json:"productId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
