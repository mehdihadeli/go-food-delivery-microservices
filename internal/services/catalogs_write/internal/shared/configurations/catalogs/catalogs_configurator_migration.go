package catalogs

import (
	"gorm.io/gorm"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/models"
)

func (ic *CatalogsServiceConfigurator) migrateCatalogs(gorm *gorm.DB) error {
	// or we could use `gorm.Migrate()`
	err := gorm.AutoMigrate(&models.Product{})
	if err != nil {
		return err
	}

	return nil
}
