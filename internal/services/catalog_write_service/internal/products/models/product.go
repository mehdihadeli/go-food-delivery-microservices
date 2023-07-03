package models

import (
	"time"

	"github.com/goccy/go-json"
	uuid "github.com/satori/go.uuid"
)

// Product model
type Product struct {
	ProductId   uuid.UUID `json:"productId"   gorm:"primaryKey"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"createdAt"` // https://gorm.io/docs/models.html#gorm-Model
	UpdatedAt   time.Time `json:"updatedAt"` // https://gorm.io/docs/models.html#gorm-Model
}

func (p *Product) String() string {
	j, _ := json.Marshal(p)
	return string(j)
}
