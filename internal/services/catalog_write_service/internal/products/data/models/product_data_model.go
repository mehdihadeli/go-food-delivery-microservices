package models

import (
	"time"

	"github.com/goccy/go-json"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// https://gorm.io/docs/conventions.html
// https://gorm.io/docs/models.html#gorm-Model

// ProductDataModel data model
type ProductDataModel struct {
	gorm.Model
	ProductId   uuid.UUID `json:"productId"   gorm:"primaryKey"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"createdAt"   gorm:"default:current_timestamp"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// https://gorm.io/docs/conventions.html#TableName
// TableName overrides the table name used by ProductDataModel to `products`
func (p *ProductDataModel) TableName() string {
	return "products"
}

func (p *ProductDataModel) String() string {
	j, _ := json.Marshal(p)

	return string(j)
}
