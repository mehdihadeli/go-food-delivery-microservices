package uow

import (
	"context"
	"emperror.dev/errors"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/constants"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger/defaultLogger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test/containers/testcontainer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/testfixture"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/data/repositories"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/models"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"testing"
	"time"
)

func Test_Catalogs_Unit_Of_Work(t *testing.T) {
	ctx := context.Background()
	gormDB, err := testcontainer.NewGormTestContainers().Start(ctx, t)
	require.NoError(t, err)

	err = seedAndMigration(ctx, gormDB)
	require.NoError(t, err)
	cfg, err := config.InitConfig(constants.Test)
	require.NoError(t, err)

	productRepository := repositories.NewPostgresProductRepository(defaultLogger.Logger, cfg, gormDB)
	uow := NewCatalogsUnitOfWork(defaultLogger.Logger, gormDB, productRepository)

	t.Run("Rollback on error", func(t *testing.T) {
		err = uow.SaveWithTx(ctx, func() error {
			_, err := uow.Products().CreateProduct(ctx, &models.Product{
				Name:        gofakeit.Name(),
				Description: gofakeit.AdjectiveDescriptive(),
				ProductId:   uuid.NewV4(),
				Price:       gofakeit.Price(100, 1000),
				CreatedAt:   time.Now(),
			})
			require.NoError(t, err)

			return errors.New("error rollback")
		})

		assert.ErrorContains(t, err, "error rollback")
		products, err := uow.Products().GetAllProducts(ctx, utils.NewListQuery(10, 1))
		require.NoError(t, err)

		assert.Equal(t, 2, len(products.Items))
	})

	t.Run("Rollback on panic", func(t *testing.T) {
		err = uow.SaveWithTx(ctx, func() error {
			_, err := uow.Products().CreateProduct(ctx, &models.Product{
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
		products, err := uow.Products().GetAllProducts(ctx, utils.NewListQuery(10, 1))
		require.NoError(t, err)

		assert.Equal(t, 2, len(products.Items))
	})

	t.Run("Rollback on context canceled", func(t *testing.T) {
		cancelCtx, cancel := context.WithCancel(ctx)

		err = uow.SaveWithTx(cancelCtx, func() error {
			_, err := uow.Products().CreateProduct(ctx, &models.Product{
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

		products, err := uow.Products().GetAllProducts(ctx, utils.NewListQuery(10, 1))
		require.NoError(t, err)

		assert.Equal(t, 2, len(products.Items))
	})

	t.Run("Commit on success", func(t *testing.T) {
		err = uow.SaveWithTx(ctx, func() error {
			_, err := uow.Products().CreateProduct(ctx, &models.Product{
				Name:        gofakeit.Name(),
				Description: gofakeit.AdjectiveDescriptive(),
				ProductId:   uuid.NewV4(),
				Price:       gofakeit.Price(100, 1000),
				CreatedAt:   time.Now(),
			})
			return err
		})
		require.NoError(t, err)
		products, err := uow.Products().GetAllProducts(ctx, utils.NewListQuery(10, 1))
		require.NoError(t, err)

		assert.Equal(t, 3, len(products.Items))
	})
}

func seedAndMigration(ctx context.Context, gormDB *gorm.DB) error {
	// migration
	err := gormDB.AutoMigrate(models.Product{})
	if err != nil {
		return err
	}

	db, err := gormDB.DB()
	if err != nil {
		return err
	}

	//https://github.com/go-testfixtures/testfixtures#templating
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

	return err
}
