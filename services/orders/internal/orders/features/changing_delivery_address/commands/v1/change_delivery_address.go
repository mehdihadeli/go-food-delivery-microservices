package v1

import (
	domainExceptions "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/exceptions/domain"
	uuid "github.com/satori/go.uuid"
)

type ChangeDeliveryAddressCommand struct {
	OrderId         uuid.UUID `validate:"required"`
	DeliveryAddress string    `validate:"required"`
}

func NewChangeDeliveryAddressCommand(orderId uuid.UUID, deliveryAddress string) (*ChangeDeliveryAddressCommand, error) {
	if deliveryAddress == "" {
		return nil, domainExceptions.ErrInvalidDeliveryAddress
	}

	return &ChangeDeliveryAddressCommand{OrderId: orderId, DeliveryAddress: deliveryAddress}, nil
}
