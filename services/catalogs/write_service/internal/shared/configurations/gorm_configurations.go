package configurations

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/gorm_postgres"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/models"
	"gorm.io/gorm"
)

func (ic *infrastructureConfigurator) configGorm() (*gorm.DB, error) {
	gorm, err := gorm_postgres.NewGorm(ic.cfg.GormPostgres)
	if err != nil {
		return nil, err
	}

	err = gorm.AutoMigrate(&models.Product{})
	if err != nil {
		return nil, err
	}

	return gorm, nil
}
