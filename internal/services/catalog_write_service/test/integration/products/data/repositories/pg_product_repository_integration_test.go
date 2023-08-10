//go:build integration
// +build integration

package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/models"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/shared/test_fixtures/integration"
)

// https://brunoscheufler.com/blog/2020-04-12-building-go-test-suites-using-testify

// Define the custom testify suite
type productPostgresRepositoryTestSuite struct {
	*integration.IntegrationTestSharedFixture
}

func TestProductPostgresRepositoryTestSuite(t *testing.T) {
	suite.Run(
		t,
		&productPostgresRepositoryTestSuite{
			IntegrationTestSharedFixture: integration.NewIntegrationTestSharedFixture(t),
		},
	)
}

func (p *productPostgresRepositoryTestSuite) Test_Create_Product_Should_Create_New_Product_In_DB() {
	ctx := context.Background()

	product := &models.Product{
		Name:        gofakeit.Name(),
		Description: gofakeit.AdjectiveDescriptive(),
		ProductId:   uuid.NewV4(),
		Price:       gofakeit.Price(100, 1000),
		CreatedAt:   time.Now(),
	}

	createdProduct, err := p.ProductRepository.CreateProduct(ctx, product)
	require.NoError(p.T(), err)

	product, err = p.ProductRepository.GetProductById(ctx, createdProduct.ProductId)
	p.NoError(err)

	p.NotNil(p)
	p.Equal(product.ProductId, createdProduct.ProductId)
}

func (p *productPostgresRepositoryTestSuite) Test_Update_Product_Should_Update_Existing_Product_In_DB() {
	ctx := context.Background()

	id := p.Items[0].ProductId
	existingProduct, err := p.ProductRepository.GetProductById(ctx, id)
	p.Require().NoError(err)
	p.Require().NotNil(existingProduct)

	existingProduct.Name = "test_update_product"
	_, err = p.ProductRepository.UpdateProduct(ctx, existingProduct)
	p.Require().NoError(err)

	updatedProduct, err := p.ProductRepository.GetProductById(ctx, id)
	p.Equal(updatedProduct.Name, "test_update_product")
}

func (p *productPostgresRepositoryTestSuite) Test_Delete_Product_Should_Delete_Existing_Product_In_DB() {
	ctx := context.Background()

	id := p.Items[0].ProductId

	err := p.ProductRepository.DeleteProductByID(ctx, id)
	p.Require().NoError(err)

	product, err := p.ProductRepository.GetProductById(ctx, id)

	p.Error(err)
	p.True(customErrors.IsNotFoundError(err))
	p.Nil(product)
}

func (p *productPostgresRepositoryTestSuite) Test_Get_Product() {
	ctx := context.Background()
	id := p.Items[0].ProductId

	p.Run("Should_Return_NotFound_Error_When_Item_DoesNot_Exists", func() {
		// with subset test a new t will create for subset test
		res, err := p.ProductRepository.GetProductById(ctx, uuid.NewV4())

		p.Error(err)
		p.True(customErrors.IsNotFoundError(err))
		p.Nil(res)
	})

	p.Run("Should_Get_Existing_Product_From_DB", func() {
		res, err := p.ProductRepository.GetProductById(ctx, id)
		p.Require().NoError(err)

		p.NotNil(res)
		p.Equal(res.ProductId, id)
	})

	p.Run("Should_Get_All_Existing_Products_From_DB", func() {
		res, err := p.ProductRepository.GetAllProducts(ctx, utils.NewListQuery(10, 1))
		p.Require().NoError(err)

		p.Equal(2, len(res.Items))
	})
}
