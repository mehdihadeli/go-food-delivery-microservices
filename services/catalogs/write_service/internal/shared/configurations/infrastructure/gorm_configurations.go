package infrastructure

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/gormPostgres"
)

func (ic *infrastructureConfigurator) configGorm() (*gormPostgres.Gorm, error) {
	gorm, err := gormPostgres.NewGorm(ic.cfg.GormPostgres)
	if err != nil {
		return nil, err
	}

	return gorm, nil
}
