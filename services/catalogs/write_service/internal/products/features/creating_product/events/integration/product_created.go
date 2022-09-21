package integration

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
	uuid "github.com/satori/go.uuid"
	"time"
)

type ProductCreated struct {
	*types.Message
	ProductId   string    `json:"productId,omitempty"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	Price       float64   `json:"price,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
}

func NewProductCreated(productId string, name string, description string, price float64, createdAt time.Time) *ProductCreated {
	return &ProductCreated{ProductId: productId, Name: name, Description: description, Price: price, CreatedAt: createdAt, Message: types.NewMessage(uuid.NewV4().String())}
}
