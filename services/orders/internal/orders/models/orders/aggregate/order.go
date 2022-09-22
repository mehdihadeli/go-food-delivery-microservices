package aggregate

//https://www.eventstore.com/blog/what-is-event-sourcing

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/domain"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/models"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mapper"
	typeMapper "github.com/mehdihadeli/store-golang-microservice-sample/pkg/reflection/type_mappper"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/serializer/jsonSerializer"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/dtos"
	domainExceptions "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/exceptions/domain"
	creatingOrderEvents "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/creating_order/events/domain/v1"
	updatingShoppingCardEvents "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/updating_shopping_card/events/v1"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/models/orders/value_objects"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Order struct {
	*models.EventSourcedAggregateRoot
	shopItems       []*value_objects.ShopItem
	accountEmail    string
	deliveryAddress string
	cancelReason    string
	totalPrice      float64
	deliveredTime   time.Time
	paid            bool
	submitted       bool
	completed       bool
	canceled        bool
	paymentId       uuid.UUID
	createdAt       time.Time
	updatedAt       time.Time
}

func (o *Order) NewEmptyAggregate() {
	//http://arch-stable.blogspot.com/2012/05/golang-call-inherited-constructor.html
	base := models.NewEventSourcedAggregateRoot(typeMapper.GetFullTypeName(o), o.When)
	o.EventSourcedAggregateRoot = base
}

func NewOrder(id uuid.UUID, shopItems []*value_objects.ShopItem, accountEmail, deliveryAddress string, deliveredTime time.Time, createdAt time.Time) (*Order, error) {
	order := &Order{}
	order.NewEmptyAggregate()
	order.SetId(id)

	if shopItems == nil || len(shopItems) == 0 {
		return nil, domainExceptions.NewOrderShopItemsRequiredError("[Order_NewOrder] order items is required")
	}

	itemsDto, err := mapper.Map[[]*dtos.ShopItemDto](shopItems)
	if err != nil {
		return nil, customErrors.NewDomainErrorWrap(err, "[Order_NewOrder.Map] error in the mapping []ShopItems to []ShopItemsDto")
	}

	event, err := creatingOrderEvents.NewOrderCreatedEventV1(id, itemsDto, accountEmail, deliveryAddress, deliveredTime, createdAt)
	if err != nil {
		return nil, customErrors.NewDomainErrorWrap(err, "[Order_NewOrder.NewOrderCreatedEventV1] error in creating order created event")
	}

	err = order.Apply(event, true)
	if err != nil {
		return nil, customErrors.NewDomainErrorWrap(err, "[Order_NewOrder.Apply] error in applying created event")
	}

	return order, nil
}

func (o *Order) UpdateShoppingCard(shopItems []*value_objects.ShopItem) error {

	event, err := updatingShoppingCardEvents.NewShoppingCartUpdatedV1(shopItems)
	if err != nil {
		return err
	}

	err = o.Apply(event, true)
	if err != nil {
		return err
	}

	return nil
}

func (o *Order) When(event domain.IDomainEvent) error {
	switch evt := event.(type) {

	case *creatingOrderEvents.OrderCreatedV1:
		return o.onOrderCreated(evt)

	default:
		return errors.InvalidEventTypeError
	}
}

func (o *Order) onOrderCreated(evt *creatingOrderEvents.OrderCreatedV1) error {
	items, err := mapper.Map[[]*value_objects.ShopItem](evt.ShopItems)
	if err != nil {
		return err
	}

	o.accountEmail = evt.AccountEmail
	o.shopItems = items
	o.deliveryAddress = evt.DeliveryAddress
	o.deliveredTime = evt.DeliveredTime
	o.createdAt = evt.CreatedAt
	o.SetId(evt.GetAggregateId()) // o.SetId(evt.Id)

	return nil
}

func (o *Order) ShopItems() []*value_objects.ShopItem {
	return o.shopItems
}

func (o *Order) PaymentId() uuid.UUID {
	return o.paymentId
}

func (o *Order) AccountEmail() string {
	return o.accountEmail
}

func (o *Order) DeliveryAddress() string {
	return o.deliveryAddress
}

func (o *Order) DeliveredTime() time.Time {
	return o.deliveredTime
}

func (o *Order) CreatedAt() time.Time {
	return o.createdAt
}

func (o *Order) TotalPrice() float64 {
	return getShopItemsTotalPrice(o.shopItems)
}

func (o *Order) Paid() bool {
	return o.paid
}

func (o *Order) Submitted() bool {
	return o.submitted
}

func (o *Order) Completed() bool {
	return o.completed
}

func (o *Order) Canceled() bool {
	return o.canceled
}

func (o *Order) CancelReason() string {
	return o.cancelReason
}

func (o *Order) String() string {
	return jsonSerializer.PrettyPrint(o)
}

func getShopItemsTotalPrice(shopItems []*value_objects.ShopItem) float64 {
	var totalPrice float64 = 0
	for _, item := range shopItems {
		totalPrice += item.Price() * float64(item.Quantity())
	}

	return totalPrice
}
