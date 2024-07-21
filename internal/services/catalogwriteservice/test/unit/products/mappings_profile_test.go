//go:build unit
// +build unit

package products

import (
	"testing"
	"time"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/mapper"
	dtoV1 "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/dtos/v1"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/models"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/shared/testfixtures/unittest"

	"github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
)

type mappingProfileUnitTests struct {
	*unittest.UnitTestSharedFixture
}

func TestMappingProfileUnit(t *testing.T) {
	suite.Run(
		t,
		&mappingProfileUnitTests{UnitTestSharedFixture: unittest.NewUnitTestSharedFixture(t)},
	)
}

func (m *mappingProfileUnitTests) Test_Mappings() {
	productModel := &models.Product{
		Id:          uuid.NewV4(),
		Name:        gofakeit.Name(),
		CreatedAt:   time.Now(),
		Description: gofakeit.EmojiDescription(),
		Price:       gofakeit.Price(100, 1000),
	}

	productDto := &dtoV1.ProductDto{
		Id:          uuid.NewV4(),
		Name:        gofakeit.Name(),
		CreatedAt:   time.Now(),
		Description: gofakeit.EmojiDescription(),
		Price:       gofakeit.Price(100, 1000),
	}

	m.Run("Should_Map_Product_To_ProductDto", func() {
		d, err := mapper.Map[*dtoV1.ProductDto](productModel)
		m.Require().NoError(err)
		m.Equal(productModel.Id, d.Id)
		m.Equal(productModel.Name, d.Name)
	})

	m.Run("Should_Map_Nil_Product_To_ProductDto", func() {
		d, err := mapper.Map[*dtoV1.ProductDto](*new(models.Product))
		m.Require().NoError(err)
		m.Nil(d)
	})

	m.Run("Should_Map_ProductDto_To_Product", func() {
		d, err := mapper.Map[*models.Product](productDto)
		m.Require().NoError(err)
		m.Equal(productDto.Id, d.Id)
		m.Equal(productDto.Name, d.Name)
	})

	m.Run("Should_Map_Nil_ProductDto_To_Product", func() {
		d, err := mapper.Map[*models.Product](*new(dtoV1.ProductDto))
		m.Require().NoError(err)
		m.Nil(d)
	})
}
