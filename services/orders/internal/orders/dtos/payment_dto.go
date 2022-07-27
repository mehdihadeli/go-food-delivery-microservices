package dtos

import "time"

type PaymentDto struct {
	PaymentId string    `json:"paymentId"`
	OrderId   string    `json:"OrderId"`
	Timestamp time.Time `json:"timestamp"`
}
