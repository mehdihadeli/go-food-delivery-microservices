package dtos

import (
	dtosV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/dtos/v1"

	customTypes "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/custom_types"
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
