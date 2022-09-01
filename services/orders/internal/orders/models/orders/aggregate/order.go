package aggregate

//https://www.eventstore.com/blog/what-is-event-sourcing

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/domain"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es"
	httpErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mapper"
	typeMapper "github.com/mehdihadeli/store-golang-microservice-sample/pkg/reflection/type_mappper"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/serializer/jsonSerializer"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/dtos"
	domainExceptions "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/exceptions/domain"
	changingDeliveryAddressEvents "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/changing_delivery_address/events/v1"
	creatingOrderEvents "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/creating_order/events/v1"
	updatingShoppingCardEvents "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/updating_shopping_card/events/v1"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/models/orders/entities"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/models/orders/value_objects"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Order struct {
	*es.EventSourcedAggregateRoot
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
	payment         *entities.Payment
	createdAt       time.Time
	updatedAt       time.Time
}

func (o *Order) NewEmptyAggregate() {
	//http://arch-stable.blogspot.com/2012/05/golang-call-inherited-constructor.html
	base := es.NewEventSourcedAggregateRoot(typeMapper.GetTypeName(o), o.When)
	o.EventSourcedAggregateRoot = base
}

func NewOrder(id uuid.UUID, shopItems []*value_objects.ShopItem, accountEmail, deliveryAddress string, deliveredTime time.Time, createdAt time.Time) (*Order, error) {
	order := &Order{}
	order.NewEmptyAggregate()
	order.SetId(id)

	return nil, domainExceptions.ErrOrderShopItemsIsRequired

	if shopItems == nil || len(shopItems) == 0 {
		return nil, domainExceptions.ErrOrderShopItemsIsRequired
	}

	itemsDto, err := mapper.Map[[]*dtos.ShopItemDto](shopItems)
	if err != nil {
		return nil, httpErrors.NewDomainErrorWrap(err, "(NewOrder): mapping shop items to dto")
	}

	event, err := creatingOrderEvents.NewOrderCreatedEventV1(id, itemsDto, accountEmail, deliveryAddress, deliveredTime, createdAt)

	if err != nil {
		return nil, httpErrors.NewDomainErrorWrap(err, "(NewOrder): error in creating order created event")
	}

	err = order.Apply(event, true)
	if err != nil {
		return nil, httpErrors.NewDomainErrorWrap(err, "(NewOrder): error in applying created event")
	}

	return order, nil
}

func (o *Order) ChangeDeliveryAddress(address string) error {

	event, err := changingDeliveryAddressEvents.NewDeliveryAddressChangedEventV1(address)
	if err != nil {
		return err
	}

	err = o.Apply(event, true)
	if err != nil {
		return err
	}

	return nil
}

func (o *Order) UpdateShoppingCard(shopItems []*value_objects.ShopItem) error {

	event, err := updatingShoppingCardEvents.NewShoppingCartUpdatedEventV1(shopItems)
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

	case *creatingOrderEvents.OrderCreatedEventV1:
		return o.onOrderCreated(evt)
	//case payingOrderEvents.OrderPaid:
	//	return o.onOrderPaid(evt)
	//case submittingOrderEvents.OrderSubmitted:
	//	return o.onOrderSubmitted(evt)
	//case completingOrderEvents.OrderCompleted:
	//	return o.onOrderCompleted(evt)
	//case cancelingOrderEvents.OrderCanceled:
	//	return o.onOrderCanceled(evt)
	//case updatingShoppingCardEvents.ShoppingCartUpdated:
	//	return o.onShoppingCartUpdated(evt)
	//case changingDeliveryAddressEvents.DeliveryAddressChanged:
	//	return o.onChangeDeliveryAddress(evt)

	default:
		return es.ErrInvalidEventType
	}
}

func (o *Order) onOrderCreated(evt *creatingOrderEvents.OrderCreatedEventV1) error {
	items, err := mapper.Map[[]*value_objects.ShopItem](evt.ShopItems)
	if err != nil {
		return err
	}

	o.accountEmail = evt.AccountEmail
	o.shopItems = items
	o.totalPrice = getShopItemsTotalPrice(items)
	o.deliveryAddress = evt.DeliveryAddress
	o.deliveredTime = evt.DeliveredTime
	o.createdAt = evt.CreatedAt
	o.SetId(evt.GetAggregateId()) // o.SetId(evt.OrderId)

	return nil
}

//func (o *Order) onOrderPaid(evt *es.Event) error {
//	var payment Payment
//	if err := evt.GetJsonData(&payment); err != nil {
//		return http_errors.Wrap(err, "GetJsonData")
//	}
//
//	o.Paid = true
//	o.Payment = payment
//
//	return nil
//}
//
//func (o *Order) onOrderSubmitted(evt *es.Event) error {
//	o.Submitted = true
//
//	return nil
//}
//
//func (o *Order) onOrderCompleted(evt *es.Event) error {
//	var eventData completingOrderEvents.OrderCompletedEvent
//	if err := evt.GetJsonData(&eventData); err != nil {
//		return http_errors.Wrap(err, "GetJsonData")
//	}
//
//	o.Completed = true
//	o.DeliveredTime = eventData.DeliveryTimestamp
//	o.Canceled = false
//
//	return nil
//}
//
//func (o *Order) onOrderCanceled(evt *es.Event) error {
//	var eventData cancelingOrderEvents.OrderCanceledEvent
//	if err := evt.GetJsonData(&eventData); err != nil {
//		return http_errors.Wrap(err, "GetJsonData")
//	}
//
//	o.Canceled = true
//	o.Completed = false
//	o.CancelReason = eventData.CancelReason
//
//	return nil
//}
//
//func (o *Order) onShoppingCartUpdated(evt *es.Event) error {
//	var eventData updatingShoppingCardEvents.ShoppingCartUpdatedEvent
//	if err := evt.GetJsonData(&eventData); err != nil {
//		return http_errors.Wrap(err, "GetJsonData")
//	}
//
//	o.ShopItems = eventData.ShopItems
//	o.TotalPrice = getShopItemsTotalPrice(eventData.ShopItems)
//
//	return nil
//}
//
//func (o *Order) onChangeDeliveryAddress(evt *es.Event) error {
//	var eventData changingDeliveryAddressEvents.DeliveryAddressChangedEvent
//	if err := evt.GetJsonData(&eventData); err != nil {
//		return http_errors.Wrap(err, "GetJsonData")
//	}
//
//	o.DeliveryAddress = eventData.DeliveryAddress
//
//	return nil
//}

func (o *Order) ShopItems() []*value_objects.ShopItem {
	return o.shopItems
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
	return o.totalPrice
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
