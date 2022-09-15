package v1

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/domain"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/models/orders/value_objects"
)

type ShoppingCartUpdatedV1 struct {
	*domain.DomainEvent
	ShopItems []*value_objects.ShopItem `json:"shopItems" bson:"shopItems,omitempty"`
}

func NewShoppingCartUpdatedV1(shopItems []*value_objects.ShopItem) (*ShoppingCartUpdatedV1, error) {
	//if shopItems == nil {
	//	return nil, domainExceptions.ErrOrderShopItemsIsRequired
	//}

	eventData := ShoppingCartUpdatedV1{ShopItems: shopItems}

	return &eventData, nil
}
