package v1

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/dtos"
	uuid "github.com/satori/go.uuid"
	"time"
)

//https://echo.labstack.com/guide/request/
//https://github.com/go-playground/validator

type CreateOrder struct {
	OrderID         uuid.UUID           `validate:"required"`
	ShopItems       []*dtos.ShopItemDto `validate:"required"`
	AccountEmail    string              `validate:"required,email"`
	DeliveryAddress string              `validate:"required"`
	DeliveryTime    time.Time           `validate:"required"`
	CreatedAt       time.Time           `validate:"required"`
}

func NewCreateOrder(shopItems []*dtos.ShopItemDto, accountEmail, deliveryAddress string, deliveryTime time.Time) *CreateOrder {
	return &CreateOrder{OrderID: uuid.NewV4(), ShopItems: shopItems, AccountEmail: accountEmail, DeliveryAddress: deliveryAddress, DeliveryTime: deliveryTime, CreatedAt: time.Now()}
}
