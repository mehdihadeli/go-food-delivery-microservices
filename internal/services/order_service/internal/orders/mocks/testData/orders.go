package testData

import (
	"github.com/brianvoe/gofakeit/v6"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/models/orders/read_models"
)

var Orders = []*read_models.OrderReadModel{
	{
		Id:              gofakeit.UUID(),
		OrderId:         gofakeit.UUID(),
		ShopItems:       generateShopItems(),
		AccountEmail:    gofakeit.Email(),
		DeliveryAddress: gofakeit.Address().Address,
		CancelReason:    gofakeit.Sentence(5),
		TotalPrice:      gofakeit.Float64Range(10, 100),
		DeliveredTime:   gofakeit.Date(),
		Paid:            gofakeit.Bool(),
		Submitted:       gofakeit.Bool(),
		Completed:       gofakeit.Bool(),
		Canceled:        gofakeit.Bool(),
		PaymentId:       gofakeit.UUID(),
		CreatedAt:       gofakeit.Date(),
		UpdatedAt:       gofakeit.Date(),
	},
	{
		Id:              gofakeit.UUID(),
		OrderId:         gofakeit.UUID(),
		ShopItems:       generateShopItems(),
		AccountEmail:    gofakeit.Email(),
		DeliveryAddress: gofakeit.Address().Address,
		CancelReason:    gofakeit.Sentence(5),
		TotalPrice:      gofakeit.Float64Range(10, 100),
		DeliveredTime:   gofakeit.Date(),
		Paid:            gofakeit.Bool(),
		Submitted:       gofakeit.Bool(),
		Completed:       gofakeit.Bool(),
		Canceled:        gofakeit.Bool(),
		PaymentId:       gofakeit.UUID(),
		CreatedAt:       gofakeit.Date(),
		UpdatedAt:       gofakeit.Date(),
	},
}

func generateShopItems() []*read_models.ShopItemReadModel {
	var shopItems []*read_models.ShopItemReadModel

	for i := 0; i < 3; i++ {
		shopItem := &read_models.ShopItemReadModel{
			Title:       gofakeit.Word(),
			Description: gofakeit.Sentence(3),
			Quantity:    uint64(gofakeit.UintRange(1, 100)),
			Price:       gofakeit.Float64Range(1, 50),
		}

		shopItems = append(shopItems, shopItem)
	}

	return shopItems
}
