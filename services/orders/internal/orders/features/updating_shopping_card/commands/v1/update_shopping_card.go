package v1

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/dtos"
	uuid "github.com/satori/go.uuid"
)

type UpdateShoppingCart struct {
	OrderId   uuid.UUID           `validate:"required"`
	ShopItems []*dtos.ShopItemDto `validate:"required"`
}

func NewUpdateShoppingCart(orderId uuid.UUID, shopItems []*dtos.ShopItemDto) *UpdateShoppingCart {
	return &UpdateShoppingCart{OrderId: orderId, ShopItems: shopItems}
}
