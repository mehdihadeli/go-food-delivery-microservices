package aggregate

//https://www.eventstore.com/blog/what-is-event-sourcing

import (
	"fmt"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/domain/types"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es"
	typeMapper "github.com/mehdihadeli/store-golang-microservice-sample/pkg/reflection/type_mappper"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/serializer"
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
	ShopItems       []*value_objects.ShopItem `json:"shopItems" bson:"shopItems,omitempty"`
	AccountEmail    string                    `json:"accountEmail" bson:"accountEmail,omitempty"`
	DeliveryAddress string                    `json:"deliveryAddress" bson:"deliveryAddress,omitempty"`
	CancelReason    string                    `json:"cancelReason" bson:"cancelReason,omitempty"`
	TotalPrice      float64                   `json:"totalPrice" bson:"totalPrice,omitempty"`
	DeliveredTime   time.Time                 `json:"deliveredTime" bson:"deliveredTime,omitempty"`
	Paid            bool                      `json:"paid" bson:"paid,omitempty"`
	Submitted       bool                      `json:"submitted" bson:"submitted,omitempty"`
	Completed       bool                      `json:"completed" bson:"completed,omitempty"`
	Canceled        bool                      `json:"canceled" bson:"canceled,omitempty"`
	Payment         *entities.Payment         `json:"payment" bson:"payment,omitempty"`
	CreatedAt       time.Time                 `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt       time.Time                 `json:"updatedAt"  bson:"createdAt"`
}

func (o *Order) NewEmptyAggregate() {
	//http://arch-stable.blogspot.com/2012/05/golang-call-inherited-constructor.html
	base := es.NewEventSourcedAggregateRoot(o.When)
	base.SetType(types.AggregateType(typeMapper.GetTypeName(o)))
	o.EventSourcedAggregateRoot = base
}

func CreateNewOrder(orderId uuid.UUID, shopItems []*value_objects.ShopItem, accountEmail, deliveryAddress string, deliveredTime time.Time, createdAt time.Time) (*Order, error) {
	order := &Order{}
	order.NewEmptyAggregate()
	order.SetID(orderId)

	event, err := creatingOrderEvents.NewOrderCreatedEvent(shopItems, accountEmail, deliveryAddress, deliveredTime, createdAt)
	if err != nil {
		return nil, err
	}

	err = order.Apply(event)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (o *Order) ChangeDeliveryAddress(address string) error {

	event, err := changingDeliveryAddressEvents.NewDeliveryAddressChangedEvent(address)
	if err != nil {
		return err
	}

	err = o.Apply(event)
	if err != nil {
		return err
	}

	return nil
}

func (o *Order) UpdateShoppingCard(shopItems []*value_objects.ShopItem) error {

	event, err := updatingShoppingCardEvents.NewShoppingCartUpdatedEvent(shopItems)
	if err != nil {
		return err
	}

	err = o.Apply(event)
	if err != nil {
		return err
	}

	return nil
}

func (o *Order) When(event interface{}) error {

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

	o.AccountEmail = evt.AccountEmail
	o.ShopItems = evt.ShopItems
	o.TotalPrice = getShopItemsTotalPrice(evt.ShopItems)
	o.DeliveryAddress = evt.DeliveryAddress
	o.DeliveredTime = evt.DeliveredTime
	o.CreatedAt = evt.CreatedAt

	return nil
}

//
//func (o *Order) onOrderPaid(evt *es.Event) error {
//	var payment Payment
//	if err := evt.GetJsonData(&payment); err != nil {
//		return errors.Wrap(err, "GetJsonData")
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
//		return errors.Wrap(err, "GetJsonData")
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
//		return errors.Wrap(err, "GetJsonData")
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
//		return errors.Wrap(err, "GetJsonData")
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
//		return errors.Wrap(err, "GetJsonData")
//	}
//
//	o.DeliveryAddress = eventData.DeliveryAddress
//
//	return nil
//}

func (o *Order) String() string {

	serializer.PrettyPrint(o)

	return fmt.Sprintf("ID: {%s}, ShopItems: {%+v}, Paid: {%v}, Submitted: {%v}, "+
		"Completed: {%v}, Canceled: {%v}, CancelReason: {%s}, TotalPrice: {%v}, AccountEmail: {%s}, DeliveryAddress: {%s}, DeliveredTime: {%s}, Payment: {%s}",
		o.ID,
		o.ShopItems,
		o.Paid,
		o.Submitted,
		o.Completed,
		o.Canceled,
		o.CancelReason,
		o.TotalPrice,
		o.AccountEmail,
		o.DeliveryAddress,
		o.DeliveredTime.String(),
		o.Payment.String(),
	)
}

func getShopItemsTotalPrice(shopItems []*value_objects.ShopItem) float64 {
	var totalPrice float64 = 0
	for _, item := range shopItems {
		totalPrice += item.Price * float64(item.Quantity)
	}

	return totalPrice
}
