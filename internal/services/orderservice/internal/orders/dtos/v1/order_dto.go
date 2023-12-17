package dtosV1

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type OrderDto struct {
	Id              uuid.UUID      `json:"id"`
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
	PaymentId       uuid.UUID      `json:"paymentId"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	OriginalVersion int64          `json:"originalVersion"`
}
