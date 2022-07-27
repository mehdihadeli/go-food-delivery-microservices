package v1

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/dtos"
	uuid "github.com/satori/go.uuid"
)

type UpdateShoppingCartCommand struct {
	OrderId   uuid.UUID           `validate:"required"`
	ShopItems []*dtos.ShopItemDto `validate:"required"`
}

func NewUpdateShoppingCartCommand(orderId uuid.UUID, shopItems []*dtos.ShopItemDto) *UpdateShoppingCartCommand {
	return &UpdateShoppingCartCommand{OrderId: orderId, ShopItems: shopItems}
}
