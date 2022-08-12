package dtos

import "time"

type PaymentDto struct {
	PaymentId string    `json:"paymentId"`
	OrderId   string    `json:"orderId"`
	Timestamp time.Time `json:"timestamp"`
}
