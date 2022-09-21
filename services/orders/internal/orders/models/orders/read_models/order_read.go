package read_models

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type OrderReadModel struct {
	// we generate id ourself because auto generate mongo string id column with type _id is not an uuid
	Id              string               `json:"id" bson:"_id,omitempty"` //https://www.mongodb.com/docs/drivers/go/current/fundamentals/crud/write-operations/insert/#the-_id-field
	OrderId         string               `json:"orderId" bson:"orderId,omitempty"`
	ShopItems       []*ShopItemReadModel `json:"shopItems,omitempty" bson:"shopItems,omitempty"`
	AccountEmail    string               `json:"accountEmail,omitempty" bson:"accountEmail,omitempty"`
	DeliveryAddress string               `json:"deliveryAddress,omitempty" bson:"deliveryAddress,omitempty"`
	CancelReason    string               `json:"cancelReason,omitempty" bson:"cancelReason,omitempty"`
	TotalPrice      float64              `json:"totalPrice,omitempty" bson:"totalPrice,omitempty"`
	DeliveredTime   time.Time            `json:"deliveredTime,omitempty" bson:"deliveredTime,omitempty"`
	Paid            bool                 `json:"paid,omitempty" bson:"paid,omitempty"`
	Submitted       bool                 `json:"submitted,omitempty" bson:"submitted,omitempty"`
	Completed       bool                 `json:"completed,omitempty" bson:"completed,omitempty"`
	Canceled        bool                 `json:"canceled,omitempty" bson:"canceled,omitempty"`
	PaymentId       string               `json:"paymentId" bson:"paymentId,omitempty"`
	CreatedAt       time.Time            `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt       time.Time            `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

func NewOrderReadModel(orderId uuid.UUID, items []*ShopItemReadModel, accountEmail string, deliveryAddress string, deliveryTime time.Time) *OrderReadModel {
	return &OrderReadModel{
		Id:              uuid.NewV4().String(), // we generate id ourself because auto generate mongo string id column with type _id is not an uuid
		OrderId:         orderId.String(),
		ShopItems:       items,
		AccountEmail:    accountEmail,
		DeliveryAddress: deliveryAddress,
		TotalPrice:      getShopItemsTotalPrice(items),
		DeliveredTime:   deliveryTime,
		CreatedAt:       time.Now(),
	}
}

func getShopItemsTotalPrice(shopItems []*ShopItemReadModel) float64 {
	var totalPrice float64 = 0
	for _, item := range shopItems {
		totalPrice += item.Price * float64(item.Quantity)
	}

	return totalPrice
}
