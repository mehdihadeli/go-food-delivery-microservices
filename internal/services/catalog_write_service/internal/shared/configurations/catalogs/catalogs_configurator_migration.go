package catalogs

import (
	"context"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/migration/contracts"

	"gorm.io/gorm"
)

func (ic *CatalogsServiceConfigurator) migrateCatalogs(
	gorm *gorm.DB,
	runner contracts.PostgresMigrationRunner,
) error {
	// - for complex migration and ability to back-track to specific migration revision it is better we use `goose`, but if we want to use built-in gorm migration we can also sync gorm with `atlas` integration migration versioning for getting migration history from grom changes
	// - here I used goose for migration, with using cmd/migration file
	// https://atlasgo.io/guides/orms/gorm
	//err := gorm.AutoMigrate(&models.Product{})
	//if err != nil {
	//	return err
	//}

	// migration with Goorse
	return ic.migrateGoose(runner)
}

func (ic *CatalogsServiceConfigurator) migrateGoose(
	runner contracts.PostgresMigrationRunner,
) error {
	err := runner.Up(context.Background(), 0)

	return err
}
