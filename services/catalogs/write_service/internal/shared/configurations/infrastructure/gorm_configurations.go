package infrastructure

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/gorm_postgres"
)

func (ic *infrastructureConfigurator) configGorm() (*gorm_postgres.Gorm, error) {
	gorm, err := gorm_postgres.NewGorm(ic.cfg.GormPostgres)
	if err != nil {
		return nil, err
	}

	return gorm, nil
}
