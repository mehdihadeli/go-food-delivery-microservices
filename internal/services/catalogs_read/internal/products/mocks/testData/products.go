package testData

import (
	"time"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/products/models"

	"github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"
)

var Products = []*models.Product{
	{
		Id:          uuid.NewV4().String(),
		ProductId:   uuid.NewV4().String(),
		Name:        gofakeit.Name(),
		CreatedAt:   time.Now(),
		Description: gofakeit.AdjectiveDescriptive(),
		Price:       gofakeit.Price(100, 1000),
	},
	{
		Id:          uuid.NewV4().String(),
		ProductId:   uuid.NewV4().String(),
		Name:        gofakeit.Name(),
		CreatedAt:   time.Now(),
		Description: gofakeit.AdjectiveDescriptive(),
		Price:       gofakeit.Price(100, 1000),
	},
}
