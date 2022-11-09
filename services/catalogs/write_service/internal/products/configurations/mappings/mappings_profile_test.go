package mappings

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mapper"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/dto"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/models"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/test_fixtures/unit_test"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type mappingProfileUnitTests struct {
	*unit_test.UnitTestSharedFixture
}

func TestMappingProfileUnit(t *testing.T) {
	suite.Run(t, &mappingProfileUnitTests{UnitTestSharedFixture: unit_test.NewUnitTestSharedFixture(t)})
}

func (m *mappingProfileUnitTests) Test_Mappings() {
	productModel := &models.Product{
		ProductId:   uuid.NewV4(),
		Name:        gofakeit.Name(),
		CreatedAt:   time.Now(),
		Description: gofakeit.EmojiDescription(),
		Price:       gofakeit.Price(100, 1000),
	}

	productDto := &dto.ProductDto{
		ProductId:   uuid.NewV4(),
		Name:        gofakeit.Name(),
		CreatedAt:   time.Now(),
		Description: gofakeit.EmojiDescription(),
		Price:       gofakeit.Price(100, 1000),
	}

	m.Run("Should_Map_Product_To_ProductDto", func() {
		d, err := mapper.Map[*dto.ProductDto](productModel)
		m.Require().NoError(err)
		m.Equal(productModel.ProductId, d.ProductId)
		m.Equal(productModel.Name, d.Name)
	})

	m.Run("Should_Map_Nil_Product_To_ProductDto", func() {
		d, err := mapper.Map[*dto.ProductDto](*new(models.Product))
		m.Require().NoError(err)
		m.Nil(d)
	})

	m.Run("Should_Map_ProductDto_To_Product", func() {
		d, err := mapper.Map[*models.Product](productDto)
		m.Require().NoError(err)
		m.Equal(productDto.ProductId, d.ProductId)
		m.Equal(productDto.Name, d.Name)
	})

	m.Run("Should_Map_Nil_ProductDto_To_Product", func() {
		d, err := mapper.Map[*models.Product](*new(dto.ProductDto))
		m.Require().NoError(err)
		m.Nil(d)
	})
}
