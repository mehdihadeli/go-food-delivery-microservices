package catalogs

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/models"

	"gorm.io/gorm"
)

func (ic *CatalogsServiceConfigurator) migrateCatalogs(gorm *gorm.DB) error {
	// or we could use `gorm.Migrate()`
	err := gorm.AutoMigrate(&models.Product{})
	if err != nil {
		return err
	}

	return nil
}
