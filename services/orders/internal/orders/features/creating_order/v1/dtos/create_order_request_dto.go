package dtos

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/custom_types"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/dtos/v1"
)

// https://echo.labstack.com/guide/binding/
// https://echo.labstack.com/guide/request/
// https://github.com/go-playground/validator

// CreateOrderRequestDto validation will handle in command level
type CreateOrderRequestDto struct {
	ShopItems       []*dtosV1.ShopItemDto  `json:"shopItems"`
	AccountEmail    string                 `json:"accountEmail"`
	DeliveryAddress string                 `json:"deliveryAddress"`
	DeliveryTime    customTypes.CustomTime `json:"deliveryTime"`
}
