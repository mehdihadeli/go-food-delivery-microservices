package v1

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
	"time"
)

type ProductUpdatedV1 struct {
	*types.Message
	ProductId   string    `json:"productId,omitempty"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	Price       float64   `json:"price,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"`
}
