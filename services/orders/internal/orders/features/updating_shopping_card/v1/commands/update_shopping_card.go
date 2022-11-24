package commands

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/dtos/v1"
	uuid "github.com/satori/go.uuid"
)

type UpdateShoppingCart struct {
	OrderId   uuid.UUID             `validate:"required"`
	ShopItems []*dtosV1.ShopItemDto `validate:"required"`
}

func NewUpdateShoppingCart(orderId uuid.UUID, shopItems []*dtosV1.ShopItemDto) *UpdateShoppingCart {
	return &UpdateShoppingCart{OrderId: orderId, ShopItems: shopItems}
}
