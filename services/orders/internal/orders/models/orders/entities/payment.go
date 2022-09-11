package entities

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type Payment struct {
	paymentId uuid.UUID
	orderId   uuid.UUID
	timestamp time.Time
}

func NewPayment(paymentId uuid.UUID, orderId uuid.UUID, timestamp time.Time) *Payment {
	return &Payment{paymentId: paymentId, orderId: orderId, timestamp: timestamp}
}

func (p *Payment) PaymentId() uuid.UUID {
	return p.paymentId
}

func (p *Payment) OrderId() uuid.UUID {
	return p.orderId
}

func (p *Payment) Timestamp() time.Time {
	return p.timestamp
}
