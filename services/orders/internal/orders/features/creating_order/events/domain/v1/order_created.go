package v1

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/domain"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	typeMapper "github.com/mehdihadeli/store-golang-microservice-sample/pkg/reflection/type_mappper"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/dtos"
	domainExceptions "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/exceptions/domain"
	uuid "github.com/satori/go.uuid"
	"time"
)

type OrderCreatedV1 struct {
	*domain.DomainEvent
	OrderId         uuid.UUID           `json:"order_id"`
	ShopItems       []*dtos.ShopItemDto `json:"shopItems" bson:"shopItems,omitempty"`
	AccountEmail    string              `json:"accountEmail" bson:"accountEmail,omitempty"`
	DeliveryAddress string              `json:"deliveryAddress" bson:"deliveryAddress,omitempty"`
	CreatedAt       time.Time           `json:"createdAt" bson:"createdAt,omitempty"`
	DeliveredTime   time.Time           `json:"deliveredTime" bson:"deliveredTime,omitempty"`
}

func NewOrderCreatedEventV1(aggregateId uuid.UUID, shopItems []*dtos.ShopItemDto, accountEmail, deliveryAddress string, deliveredTime time.Time, createdAt time.Time) (*OrderCreatedV1, error) {
	if shopItems == nil || len(shopItems) == 0 {
		return nil, domainExceptions.NewOrderShopItemsRequiredError("shopItems is required")
	}

	if deliveryAddress == "" {
		return nil, domainExceptions.NewInvalidDeliveryAddressError("deliveryAddress is invalid")
	}

	if accountEmail == "" {
		return nil, domainExceptions.NewInvalidEmailAddressError("accountEmail is invalid")
	}

	if createdAt.IsZero() {
		return nil, customErrors.NewDomainError("createdAt can't be zero")
	}

	if deliveredTime.IsZero() {
		return nil, customErrors.NewDomainError("deliveredTime can't be zero")
	}

	eventData := &OrderCreatedV1{
		ShopItems:       shopItems,
		OrderId:         aggregateId,
		AccountEmail:    accountEmail,
		DeliveryAddress: deliveryAddress,
		CreatedAt:       createdAt,
		DeliveredTime:   deliveredTime,
	}

	eventData.DomainEvent = domain.NewDomainEvent(typeMapper.GetTypeName(eventData))

	return eventData, nil
}
