package v1

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/domain"
)

type DeliveryAddressChangedEventV1 struct {
	*domain.DomainEvent
	DeliveryAddress string `json:"deliveryAddress" bson:"deliveryAddress,omitempty"`
}

func NewDeliveryAddressChangedEventV1(deliveryAddress string) (*DeliveryAddressChangedEventV1, error) {
	//if deliveryAddress == "" {
	//	return nil, domainExceptions.ErrInvalidDeliveryAddress
	//}

	eventData := DeliveryAddressChangedEventV1{DeliveryAddress: deliveryAddress}

	return &eventData, nil
}
