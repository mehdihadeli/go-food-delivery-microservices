package commands

import (
	uuid "github.com/satori/go.uuid"

	dtosV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/dtos/v1"
)

type UpdateShoppingCart struct {
	OrderId   uuid.UUID             `validate:"required"`
	ShopItems []*dtosV1.ShopItemDto `validate:"required"`
}

func NewUpdateShoppingCart(orderId uuid.UUID, shopItems []*dtosV1.ShopItemDto) *UpdateShoppingCart {
	return &UpdateShoppingCart{OrderId: orderId, ShopItems: shopItems}
}
