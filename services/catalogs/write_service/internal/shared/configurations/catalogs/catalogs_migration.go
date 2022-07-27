package catalogs

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/gorm_postgres"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/models"
)

func (c *catalogsServiceConfigurator) migrateCatalogs(gorm *gorm_postgres.Gorm) error {

	// or we could use `gorm.Migrate()`
	err := gorm.DB.AutoMigrate(&models.Product{})
	if err != nil {
		return err
	}

	return nil
}
