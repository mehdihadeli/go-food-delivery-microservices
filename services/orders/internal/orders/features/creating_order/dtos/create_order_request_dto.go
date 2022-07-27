package dtos

import (
	customTypes "github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/custom_types"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/dtos"
)

type CreateOrderRequestDto struct {
	ShopItems       []*dtos.ShopItemDto    `validate:"required"`
	AccountEmail    string                 `validate:"required,email"`
	DeliveryAddress string                 `validate:"required"`
	DeliveryTime    customTypes.CustomTime `validate:"required"`
}
