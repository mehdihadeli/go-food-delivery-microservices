package catalogs

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/models"
	"gorm.io/gorm"
)

func (c *catalogsServiceConfigurator) migrateCatalogs(gorm *gorm.DB) error {
	// or we could use `gorm.Migrate()`
	err := gorm.AutoMigrate(&models.Product{})
	if err != nil {
		return err
	}

	return nil
}
