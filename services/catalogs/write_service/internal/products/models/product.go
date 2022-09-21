package models

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/serializer/jsonSerializer"
	"time"

	uuid "github.com/satori/go.uuid"
)

// Product model
type Product struct {
	ProductId   uuid.UUID `json:"productId" gorm:"primaryKey"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (p *Product) String() string {
	return jsonSerializer.PrettyPrint(p)
}
