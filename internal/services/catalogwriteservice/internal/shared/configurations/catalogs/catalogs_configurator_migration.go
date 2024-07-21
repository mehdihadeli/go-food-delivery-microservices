package catalogs

import (
	"context"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/migration/contracts"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/models"

	"gorm.io/gorm"
)

func (ic *CatalogsServiceConfigurator) migrateCatalogs(
	runner contracts.PostgresMigrationRunner,
) error {
	// - for complex migration and ability to back-track to specific migration revision it is better we use `goose`, but if we want to use built-in gorm migration we can also sync gorm with `atlas` integration migration versioning for getting migration history from grom changes
	// - here I used goose for migration, with using cmd/migration file

	// migration with Goorse
	return migrateGoose(runner)
}

func migrateGoose(
	runner contracts.PostgresMigrationRunner,
) error {
	err := runner.Up(context.Background(), 0)

	return err
}

func migrateGorm(
	db *gorm.DB,
) error {
	// https://atlasgo.io/guides/orms/gorm
	err := db.AutoMigrate(&models.Product{})
	if err != nil {
		return err
	}

	return nil
}
