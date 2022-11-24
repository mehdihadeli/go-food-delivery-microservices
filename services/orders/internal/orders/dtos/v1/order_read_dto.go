package dtosV1

import "time"

type OrderReadDto struct {
	Id              string             `json:"id"`
	OrderId         string             `json:"orderId"`
	ShopItems       []*ShopItemReadDto `json:"shopItems"`
	AccountEmail    string             `json:"accountEmail"`
	DeliveryAddress string             `json:"deliveryAddress"`
	CancelReason    string             `json:"cancelReason"`
	TotalPrice      float64            `json:"totalPrice"`
	DeliveredTime   time.Time          `json:"deliveredTime"`
	Paid            bool               `json:"paid"`
	Submitted       bool               `json:"submitted"`
	Completed       bool               `json:"completed"`
	Canceled        bool               `json:"canceled"`
	PaymentId       string             `json:"paymentId"`
	CreatedAt       time.Time          `json:"createdAt"`
	UpdatedAt       time.Time          `json:"updatedAt"`
}
