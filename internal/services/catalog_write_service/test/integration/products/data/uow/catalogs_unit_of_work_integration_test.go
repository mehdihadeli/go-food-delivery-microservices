//go:build integration
// +build integration

package uow

import (
	"context"
	"testing"
	"time"

	"emperror.dev/errors"
	"github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	data2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/contracts/data"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/models"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/shared/test_fixtures/integration"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils"
)

// https://brunoscheufler.com/blog/2020-04-12-building-go-test-suites-using-testify

// Define the custom testify suite
type catalogsUnitOfWorkTestSuite struct {
	*integration.IntegrationTestSharedFixture
}

func TestCatalogsUnitOfWorkTestSuite(t *testing.T) {
	suite.Run(
		t,
		&catalogsUnitOfWorkTestSuite{
			IntegrationTestSharedFixture: integration.NewIntegrationTestSharedFixture(t),
		},
	)
}

func (p *catalogsUnitOfWorkTestSuite) Test_Catalogs_Unit_Of_Work() {
	ctx := context.Background()

	p.T().Run("Should_Rollback_On_Error", func(t *testing.T) {
		err := p.CatalogUnitOfWorks.Do(ctx, func(catalogContext data2.CatalogContext) error {
			_, err := catalogContext.Products().CreateProduct(ctx,
				&models.Product{
					Name:        gofakeit.Name(),
					Description: gofakeit.AdjectiveDescriptive(),
					ProductId:   uuid.NewV4(),
					Price:       gofakeit.Price(100, 1000),
					CreatedAt:   time.Now(),
				})
			require.NoError(t, err)
			return errors.New("error rollback")
		})
		require.ErrorContains(t, err, "error rollback")

		products, err := p.ProductRepository.GetAllProducts(ctx, utils.NewListQuery(10, 1))
		require.NoError(t, err)

		assert.Equal(t, 2, len(products.Items))
	})

	p.T().Run("Should_Rollback_On_Panic", func(t *testing.T) {
		err := p.CatalogUnitOfWorks.Do(ctx, func(catalogContext data2.CatalogContext) error {
			_, err := catalogContext.Products().CreateProduct(ctx,
				&models.Product{
					Name:        gofakeit.Name(),
					Description: gofakeit.AdjectiveDescriptive(),
					ProductId:   uuid.NewV4(),
					Price:       gofakeit.Price(100, 1000),
					CreatedAt:   time.Now(),
				})
			require.NoError(t, err)
			panic(errors.New("panic rollback"))

			return err
		})
		require.Error(t, err)

		products, err := p.ProductRepository.GetAllProducts(ctx, utils.NewListQuery(10, 1))
		require.NoError(t, err)

		assert.Equal(t, 2, len(products.Items))
	})

	p.T().Run("Should_Rollback_On_Context_Canceled", func(t *testing.T) {
		cancelCtx, cancel := context.WithCancel(ctx)

		err := p.CatalogUnitOfWorks.Do(cancelCtx, func(catalogContext data2.CatalogContext) error {
			_, err := catalogContext.Products().CreateProduct(ctx, &models.Product{
				Name:        gofakeit.Name(),
				Description: gofakeit.AdjectiveDescriptive(),
				ProductId:   uuid.NewV4(),
				Price:       gofakeit.Price(100, 1000),
				CreatedAt:   time.Now(),
			})
			require.NoError(t, err)

			_, err = catalogContext.Products().CreateProduct(ctx, &models.Product{
				Name:        gofakeit.Name(),
				Description: gofakeit.AdjectiveDescriptive(),
				ProductId:   uuid.NewV4(),
				Price:       gofakeit.Price(100, 1000),
				CreatedAt:   time.Now(),
			})
			require.NoError(t, err)
			cancel()

			return err
		})

		products, err := p.ProductRepository.GetAllProducts(ctx, utils.NewListQuery(10, 1))
		require.NoError(t, err)

		assert.Equal(t, 2, len(products.Items))
	})

	p.T().Run("Should_Commit_On_Success", func(t *testing.T) {
		err := p.CatalogUnitOfWorks.Do(ctx, func(catalogContext data2.CatalogContext) error {
			_, err := catalogContext.Products().CreateProduct(ctx, &models.Product{
				Name:        gofakeit.Name(),
				Description: gofakeit.AdjectiveDescriptive(),
				ProductId:   uuid.NewV4(),
				Price:       gofakeit.Price(100, 1000),
				CreatedAt:   time.Now(),
			})
			return err
		})
		require.NoError(t, err)
		products, err := p.ProductRepository.GetAllProducts(ctx, utils.NewListQuery(10, 1))
		require.NoError(t, err)

		assert.Equal(t, 3, len(products.Items))
	})
}
