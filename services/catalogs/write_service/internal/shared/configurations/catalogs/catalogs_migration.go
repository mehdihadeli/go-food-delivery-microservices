package catalogs

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/gormPostgres"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/models"
)

func (c *catalogsServiceConfigurator) migrateCatalogs(gorm *gormPostgres.Gorm) error {

	// or we could use `gorm.Migrate()`
	err := gorm.DB.AutoMigrate(&models.Product{})
	if err != nil {
		return err
	}

	return nil
}
