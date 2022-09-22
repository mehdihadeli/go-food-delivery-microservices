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
	CreatedAt   time.Time `json:"createdAt"` //https://gorm.io/docs/models.html#gorm-Model
	UpdatedAt   time.Time `json:"updatedAt"` //https://gorm.io/docs/models.html#gorm-Model
}

func (p *Product) String() string {
	return jsonSerializer.PrettyPrint(p)
}
