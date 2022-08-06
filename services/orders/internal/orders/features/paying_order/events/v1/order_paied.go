package v1

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/models/orders/entities"
	uuid "github.com/satori/go.uuid"
)

type OrderPayedEventV1 struct {
	OrderID   uuid.UUID `json:"orderID" bson:"orderID,omitempty"`
	PaymentID uuid.UUID `json:"PaymentID" bson:"PaymentID,omitempty"`
}

func NewOrderPaidEventV1(payment *entities.Payment) (*OrderPayedEventV1, error) {
	event := OrderPayedEventV1{OrderID: payment.OrderId, PaymentID: payment.PaymentId}

	return &event, nil
}
