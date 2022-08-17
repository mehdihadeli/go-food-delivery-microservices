package v1

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/domain"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/models/orders/value_objects"
)

type ShoppingCartUpdatedEventV1 struct {
	*domain.DomainEvent
	ShopItems []*value_objects.ShopItem `json:"shopItems" bson:"shopItems,omitempty"`
}

func NewShoppingCartUpdatedEventV1(shopItems []*value_objects.ShopItem) (*ShoppingCartUpdatedEventV1, error) {

	//if shopItems == nil {
	//	return nil, domainExceptions.ErrOrderShopItemsIsRequired
	//}

	eventData := ShoppingCartUpdatedEventV1{ShopItems: shopItems}

	return &eventData, nil
}
