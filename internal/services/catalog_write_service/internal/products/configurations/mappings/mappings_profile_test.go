package mappings

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/mapper"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"

	dtoV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/dto/v1"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/models"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/shared/test_fixtures/unit_test"
)

type mappingProfileUnitTests struct {
	*unit_test.UnitTestSharedFixture
}

func TestMappingProfileUnit(t *testing.T) {
	suite.Run(
		t,
		&mappingProfileUnitTests{UnitTestSharedFixture: unit_test.NewUnitTestSharedFixture(t)},
	)
}

func (m *mappingProfileUnitTests) Test_Mappings() {
	productModel := &models.Product{
		ProductId:   uuid.NewV4(),
		Name:        gofakeit.Name(),
		CreatedAt:   time.Now(),
		Description: gofakeit.EmojiDescription(),
		Price:       gofakeit.Price(100, 1000),
	}

	productDto := &dtoV1.ProductDto{
		ProductId:   uuid.NewV4(),
		Name:        gofakeit.Name(),
		CreatedAt:   time.Now(),
		Description: gofakeit.EmojiDescription(),
		Price:       gofakeit.Price(100, 1000),
	}

	m.Run("Should_Map_Product_To_ProductDto", func() {
		d, err := mapper.Map[*dtoV1.ProductDto](productModel)
		m.Require().NoError(err)
		m.Equal(productModel.ProductId, d.ProductId)
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
		m.Equal(productDto.ProductId, d.ProductId)
		m.Equal(productDto.Name, d.Name)
	})

	m.Run("Should_Map_Nil_ProductDto_To_Product", func() {
		d, err := mapper.Map[*models.Product](*new(dtoV1.ProductDto))
		m.Require().NoError(err)
		m.Nil(d)
	})
}
