package createOrderCommandV1

import (
	"time"

	dtosV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/dtos/v1"

	validation "github.com/go-ozzo/ozzo-validation"
	uuid "github.com/satori/go.uuid"
)

// https://echo.labstack.com/guide/request/
// https://github.com/go-playground/validator
type CreateOrder struct {
	OrderId         uuid.UUID
	ShopItems       []*dtosV1.ShopItemDto
	AccountEmail    string
	DeliveryAddress string
	DeliveryTime    time.Time
	CreatedAt       time.Time
}

func NewCreateOrder(
	shopItems []*dtosV1.ShopItemDto,
	accountEmail, deliveryAddress string,
	deliveryTime time.Time,
) (*CreateOrder, error) {
	command := &CreateOrder{
		OrderId:         uuid.NewV4(),
		ShopItems:       shopItems,
		AccountEmail:    accountEmail,
		DeliveryAddress: deliveryAddress,
		DeliveryTime:    deliveryTime,
		CreatedAt:       time.Now(),
	}

	err := command.Validate()
	if err != nil {
		return nil, err
	}

	return command, nil
}

func (c CreateOrder) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.OrderId, validation.Required),
		validation.Field(&c.ShopItems, validation.Required),
		validation.Field(&c.AccountEmail, validation.Required),
		validation.Field(&c.DeliveryAddress, validation.Required),
		validation.Field(&c.DeliveryTime, validation.Required),
		validation.Field(&c.CreatedAt, validation.Required),
	)
}
