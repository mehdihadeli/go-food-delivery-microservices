package uow

import (
	"context"
	"testing"
	"time"

	"emperror.dev/errors"

	data2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/contracts/data"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/data/repositories"

	"github.com/brianvoe/gofakeit/v6"
	gormPostgres "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/gorm_postgres"
	defaultLogger2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/default_logger"
	gorm2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/containers/testcontainer/gorm"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/testfixture"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/models"
)

var items []*models.Product

func Test_Catalogs_Unit_Of_Work(t *testing.T) {
	ctx := context.Background()
	gormDB, err := gorm2.NewGormTestContainers().Start(ctx, t)
	require.NoError(t, err)

	err = seedAndMigration(gormDB)
	require.NoError(t, err)

	productRepository := repositories.NewPostgresProductRepository(defaultLogger2.Logger, gormDB)
	uow := NewCatalogsUnitOfWork(defaultLogger2.Logger, gormDB)

	t.Run("Should_Rollback_On_Error", func(t *testing.T) {
		err := uow.Do(ctx, func(catalogContext data2.CatalogContext) error {
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

		products, err := productRepository.GetAllProducts(ctx, utils.NewListQuery(10, 1))
		require.NoError(t, err)

		assert.Equal(t, 2, len(products.Items))
	})

	t.Run("Should_Rollback_On_Panic", func(t *testing.T) {
		err = uow.Do(ctx, func(catalogContext data2.CatalogContext) error {
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

		products, err := productRepository.GetAllProducts(ctx, utils.NewListQuery(10, 1))
		require.NoError(t, err)

		assert.Equal(t, 2, len(products.Items))
	})

	t.Run("Should_Rollback_On_Context_Canceled", func(t *testing.T) {
		cancelCtx, cancel := context.WithCancel(ctx)

		err = uow.Do(cancelCtx, func(catalogContext data2.CatalogContext) error {
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

		products, err := productRepository.GetAllProducts(ctx, utils.NewListQuery(10, 1))
		require.NoError(t, err)

		assert.Equal(t, 2, len(products.Items))
	})

	t.Run("Should_Commit_On_Success", func(t *testing.T) {
		err = uow.Do(ctx, func(catalogContext data2.CatalogContext) error {
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
		products, err := productRepository.GetAllProducts(ctx, utils.NewListQuery(10, 1))
		require.NoError(t, err)

		assert.Equal(t, 3, len(products.Items))
	})
}

func seedAndMigration(gormDB *gorm.DB) error {
	// migration
	err := gormDB.AutoMigrate(models.Product{})
	if err != nil {
		return err
	}

	db, err := gormDB.DB()
	if err != nil {
		return err
	}

	// https://github.com/go-testfixtures/testfixtures#templating
	// seed data
	err = testfixture.RunPostgresFixture(
		db,
		[]string{"db/fixtures/products"},
		map[string]interface{}{
			"Products": []struct {
				Name        string
				ProductId   uuid.UUID
				Description string
			}{
				{Name: gofakeit.Name(), Description: gofakeit.AdjectiveDescriptive(), ProductId: uuid.NewV4()},
				{Name: gofakeit.Name(), Description: gofakeit.AdjectiveDescriptive(), ProductId: uuid.NewV4()},
			},
		})
	if err != nil {
		return err
	}

	result, err := gormPostgres.Paginate[*models.Product](context.Background(), utils.NewListQuery(10, 1), gormDB)
	items = result.Items

	return err
}
