package v1

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/domain"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/models/orders/value_objects"
	"time"
)

type OrderCreatedEventV1 struct {
	*domain.DomainEvent
	ShopItems       []*value_objects.ShopItem `json:"shopItems" bson:"shopItems,omitempty"`
	AccountEmail    string                    `json:"accountEmail" bson:"accountEmail,omitempty"`
	DeliveryAddress string                    `json:"deliveryAddress" bson:"deliveryAddress,omitempty"`
	CreatedAt       time.Time                 `json:"createdAt" bson:"createdAt,omitempty"`
	DeliveredTime   time.Time                 `json:"deliveredTime" bson:"deliveredTime,omitempty"`
}

func NewOrderCreatedEventV1(shopItems []*value_objects.ShopItem, accountEmail, deliveryAddress string, deliveredTime time.Time, createdAt time.Time) (*OrderCreatedEventV1, error) {

	//if shopItems == nil {
	//	return nil, domainExceptions.ErrOrderShopItemsIsRequired
	//}
	//
	//if deliveryAddress == "" {
	//	return nil, domainExceptions.ErrInvalidDeliveryAddress
	//}
	//
	//if accountEmail == "" {
	//	return nil, domainExceptions.ErrInvalidAccountEmail
	//}
	//
	//if createdAt.IsZero() {
	//	return nil, domainExceptions.ErrInvalidTime
	//}
	//
	//if deliveredTime.IsZero() {
	//	return nil, domainExceptions.ErrInvalidTime
	//}

	eventData := &OrderCreatedEventV1{
		ShopItems:       shopItems,
		AccountEmail:    accountEmail,
		DeliveryAddress: deliveryAddress,
		CreatedAt:       createdAt,
		DeliveredTime:   deliveredTime,
	}

	return eventData, nil
}
