package entities

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Payment struct {
	paymentId uuid.UUID
	orderId   uuid.UUID
	timestamp time.Time
}

func (p *Payment) String() string {
	if p == nil {
		return "nil"
	}

	return fmt.Sprintf("PaymentID: {%s}, OrderId: {%s},  Timestamp: {%s}", p.paymentId, p, p.orderId, p.timestamp)
}
