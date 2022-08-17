package data

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/models/orders/entities"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/models/orders/value_objects"
	"time"
)

type OrderDataModel struct {
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

type OrderItemDataModel struct {
	Title       string  `json:"title" bson:"title,omitempty"`
	Description string  `json:"description" bson:"description,omitempty"`
	Quantity    uint64  `json:"quantity" bson:"quantity,omitempty"`
	Price       float64 `json:"price" bson:"price,omitempty"`
}
