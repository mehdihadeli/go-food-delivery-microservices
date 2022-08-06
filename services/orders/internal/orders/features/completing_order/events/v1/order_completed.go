package v1

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/contracts"
	domainExceptions "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/exceptions/domain"
	"time"
)

type OrderCompletedEventV1 struct {
	DeliveryTimestamp time.Time `json:"deliveryTimestamp"`
}

func NewOrderCompletedEventV1(aggregate contracts.IEventSourcedAggregateRoot, deliveryTimestamp time.Time) (*OrderCompletedEventV1, error) {
	if deliveryTimestamp.IsZero() {
		return nil, domainExceptions.ErrInvalidDeliveryTimeStamp
	}

	eventData := OrderCompletedEventV1{DeliveryTimestamp: deliveryTimestamp}

	return &eventData, nil
}
