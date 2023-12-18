package datamodels

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
	Id          uuid.UUID `gorm:"primaryKey"`
	Name        string
	Description string
	Price       float64
	CreatedAt   time.Time `gorm:"default:current_timestamp"`
	UpdatedAt   time.Time
	// for soft delete - https://gorm.io/docs/delete.html#Soft-Delete
	gorm.DeletedAt
}

// TableName overrides the table name used by ProductDataModel to `products` - https://gorm.io/docs/conventions.html#TableName
func (p *ProductDataModel) TableName() string {
	return "products"
}

func (p *ProductDataModel) String() string {
	j, _ := json.Marshal(p)

	return string(j)
}
