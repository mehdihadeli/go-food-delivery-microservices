package catalogs

import (
	"fmt"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/postgres"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/models"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

func (c *catalogsServiceConfigurator) migrateCatalogs(gorm *gorm.DB) error {

	err := createDB(c.cfg.Postgresql)
	if err != nil {
		return err
	}

	err = gorm.AutoMigrate(&models.Product{})
	if err != nil {
		return err
	}

	return nil
}

func createDB(cfg *postgres.Config) error {

	db, err := postgres.NewPgxConn(&postgres.Config{Host: cfg.Host, Port: cfg.Port, SSLMode: cfg.SSLMode, User: cfg.User, Password: cfg.Password})
	if err != nil {
		return err
	}

	var exists int
	rows, err := db.Query(context.Background(), fmt.Sprintf("SELECT 1 FROM  pg_catalog.pg_database WHERE datname='%s'", cfg.DBName))
	if err != nil {
		return err
	}

	if rows.Next() {
		err = rows.Scan(&exists)
		if err != nil {
			return err
		}
	}

	if exists == 1 {
		return nil
	}

	_, err = db.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %s", cfg.DBName))
	if err != nil {
		return err
	}

	defer db.Close()

	return nil
}
