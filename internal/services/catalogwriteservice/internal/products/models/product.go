package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// Product model
type Product struct {
	Id          uuid.UUID
	Name        string
	Description string
	Price       float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
