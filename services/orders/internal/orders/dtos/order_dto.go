package dtos

import (
	"time"
)

type OrderDto struct {
	Id              string         `json:"id"`
	ShopItems       []*ShopItemDto `json:"shopItems"`
	AccountEmail    string         `json:"accountEmail"`
	DeliveryAddress string         `json:"deliveryAddress"`
	CancelReason    string         `json:"cancelReason"`
	TotalPrice      float64        `json:"totalPrice"`
	DeliveredTime   time.Time      `json:"deliveredTime"`
	Paid            bool           `json:"paid"`
	Submitted       bool           `json:"submitted"`
	Completed       bool           `json:"completed"`
	Canceled        bool           `json:"canceled"`
	Payment         *PaymentDto    `json:"payment"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
}
