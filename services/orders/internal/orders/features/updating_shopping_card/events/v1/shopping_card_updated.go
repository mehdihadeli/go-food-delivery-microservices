package v1

import (
	domainExceptions "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/exceptions/domain"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/models/orders/value_objects"
)

type ShoppingCartUpdatedEventV1 struct {
	ShopItems []*value_objects.ShopItem `json:"shopItems" bson:"shopItems,omitempty"`
}

func NewShoppingCartUpdatedEvent(shopItems []*value_objects.ShopItem) (*ShoppingCartUpdatedEventV1, error) {

	if shopItems == nil {
		return nil, domainExceptions.ErrOrderShopItemsIsRequired
	}

	eventData := ShoppingCartUpdatedEventV1{ShopItems: shopItems}

	return &eventData, nil
}
