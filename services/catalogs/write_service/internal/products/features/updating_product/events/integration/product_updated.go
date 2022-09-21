package integration

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
	uuid "github.com/satori/go.uuid"
	"time"
)

type ProductUpdated struct {
	*types.Message
	ProductId   string    `json:"productId,omitempty"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	Price       float64   `json:"price,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"`
}

func NewProductUpdated(productID string, name string, description string, price float64) *ProductUpdated {
	return &ProductUpdated{Message: types.NewMessage(uuid.NewV4().String()), ProductId: productID, Name: name, Description: description, Price: price, UpdatedAt: time.Now()}
}
