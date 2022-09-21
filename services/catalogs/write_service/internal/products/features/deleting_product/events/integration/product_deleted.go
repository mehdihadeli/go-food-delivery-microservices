package integration

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
	uuid "github.com/satori/go.uuid"
)

type ProductDeleted struct {
	*types.Message
	ProductId string `json:"productId,omitempty"`
}

func NewProductDeleted(productId string) *ProductDeleted {
	return &ProductDeleted{ProductId: productId, Message: types.NewMessage(uuid.NewV4().String())}
}
