package v1

import domainExceptions "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/exceptions/domain"

type DeliveryAddressChangedEventV1 struct {
	DeliveryAddress string `json:"deliveryAddress" bson:"deliveryAddress,omitempty"`
}

func NewDeliveryAddressChangedEventV1(deliveryAddress string) (*DeliveryAddressChangedEventV1, error) {
	if deliveryAddress == "" {
		return nil, domainExceptions.ErrInvalidDeliveryAddress
	}

	eventData := DeliveryAddressChangedEventV1{DeliveryAddress: deliveryAddress}

	return &eventData, nil
}
