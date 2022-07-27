package entities

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Payment struct {
	PaymentId uuid.UUID `json:"paymentId" bson:"paymentId,omitempty" validate:"required"`
	OrderId   uuid.UUID `json:"orderId" bson:"orderId,omitempty" validate:"required"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp,omitempty" validate:"required"`
}

func (p *Payment) String() string {
	if p == nil {
		return "nil"
	}

	return fmt.Sprintf("PaymentID: {%s}, OrderId: {%s},  Timestamp: {%s}", p.PaymentId, p, p.OrderId, p.Timestamp)
}
