package dtos

import (
	"time"
)

type OrderDto struct {
	id              string         `json:"id"`
	shopItems       []*ShopItemDto `json:"shopItems"`
	accountEmail    string         `json:"accountEmail"`
	aeliveryAddress string         `json:"deliveryAddress"`
	cancelReason    string         `json:"cancelReason"`
	totalPrice      float64        `json:"totalPrice"`
	deliveredTime   time.Time      `json:"deliveredTime"`
	paid            bool           `json:"paid"`
	submitted       bool           `json:"submitted"`
	completed       bool           `json:"completed"`
	canceled        bool           `json:"canceled"`
	payment         *PaymentDto    `json:"payment"`
	createdAt       time.Time      `json:"createdAt"`
	updatedAt       time.Time      `json:"updatedAt"`
}
