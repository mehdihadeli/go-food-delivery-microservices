package catalogs

import (
	"time"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/testfixture"
	datamodel "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/data/datamodels"

	"emperror.dev/errors"
	"github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func (ic *CatalogsServiceConfigurator) seedCatalogs(
	db *gorm.DB,
) error {
	err := seedDataManually(db)
	if err != nil {
		return err
	}

	return nil
}

func seedDataManually(gormDB *gorm.DB) error {
	var count int64

	// https://gorm.io/docs/advanced_query.html#Count
	gormDB.Model(&datamodel.ProductDataModel{}).Count(&count)
	if count > 0 {
		return nil
	}

	products := []*datamodel.ProductDataModel{
		{
			Id:          uuid.NewV4(),
			Name:        gofakeit.Name(),
			CreatedAt:   time.Now(),
			Description: gofakeit.AdjectiveDescriptive(),
			Price:       gofakeit.Price(100, 1000),
		},
		{
			Id:          uuid.NewV4(),
			Name:        gofakeit.Name(),
			CreatedAt:   time.Now(),
			Description: gofakeit.AdjectiveDescriptive(),
			Price:       gofakeit.Price(100, 1000),
		},
	}

	err := gormDB.CreateInBatches(products, len(products)).Error
	if err != nil {
		return errors.Wrap(err, "error in seed database")
	}

	return nil
}

func seedDataWithFixture(gormDB *gorm.DB) error {
	var count int64

	// https://gorm.io/docs/advanced_query.html#Count
	gormDB.Model(&datamodel.ProductDataModel{}).Count(&count)
	if count > 0 {
		return nil
	}

	db, err := gormDB.DB()
	if err != nil {
		return errors.WrapIf(err, "error in seed database")
	}

	// https://github.com/go-testfixtures/testfixtures#templating
	// seed data
	var data []struct {
		Name        string
		ProductId   uuid.UUID
		Description string
	}

	f := []struct {
		Name        string
		ProductId   uuid.UUID
		Description string
	}{
		{gofakeit.Name(), uuid.NewV4(), gofakeit.AdjectiveDescriptive()},
		{gofakeit.Name(), uuid.NewV4(), gofakeit.AdjectiveDescriptive()},
	}

	data = append(data, f...)

	err = testfixture.RunPostgresFixture(
		db,
		[]string{"db/fixtures/products"},
		map[string]interface{}{
			"Products": data,
		})
	if err != nil {
		return errors.WrapIf(err, "error in seed database")
	}

	return nil
}
