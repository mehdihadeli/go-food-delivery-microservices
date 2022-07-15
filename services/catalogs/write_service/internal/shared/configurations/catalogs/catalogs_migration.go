package catalogs

import (
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zapadapter"
	postgres "github.com/mehdihadeli/store-golang-microservice-sample/pkg/postgres_pgx"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/models"
	"go.uber.org/zap"
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

	db, err := postgres.NewPgxPoolConn(&postgres.Config{Host: cfg.Host, Port: cfg.Port, SSLMode: cfg.SSLMode, User: cfg.User, Password: cfg.Password}, zapadapter.NewLogger(zap.L()), pgx.LogLevelInfo)
	if err != nil {
		return err
	}

	var exists int
	rows, err := db.ConnPool.Query(context.Background(), fmt.Sprintf("SELECT 1 FROM  pg_catalog.pg_database WHERE datname='%s'", cfg.DBName))
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

	_, err = db.ConnPool.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %s", cfg.DBName))
	if err != nil {
		return err
	}

	defer db.Close()

	return nil
}
