package v1

import (
	domainExceptions "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/exceptions/domain"
	uuid "github.com/satori/go.uuid"
)

type OrderSubmittedV1 struct {
	OrderID uuid.UUID `json:"orderID" bson:"orderID,omitempty"`
}

func NewSubmitOrderEvent(orderID uuid.UUID) (*OrderSubmittedV1, error) {
	if orderID == uuid.Nil {
		return nil, domainExceptions.ErrInvalidOrderID
	}

	event := OrderSubmittedV1{OrderID: orderID}

	return &event, nil
}
