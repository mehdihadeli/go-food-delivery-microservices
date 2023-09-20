package catalogs

import (
	"gorm.io/gorm"
)

func (ic *CatalogsServiceConfigurator) migrateCatalogs(gorm *gorm.DB) error {
	// - for complex migration and ability to back-track to specific migration revision it is better we use `goose`, but if we want to use built-in gorm migration we can also sync gorm with `atlas` integration migration versioning for getting migration history from grom changes
	// - here I used goose for migration, with using cmd/migration file
	// https://atlasgo.io/guides/orms/gorm

	//err := gorm.AutoMigrate(&models.Product{})
	//if err != nil {
	//	return err
	//}

	return nil
}
