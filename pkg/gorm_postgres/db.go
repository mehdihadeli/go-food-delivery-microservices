package gorm_postgres

import (
	"context"
	"fmt"
	postgres "github.com/mehdihadeli/store-golang-microservice-sample/pkg/postgres"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	gorm_postgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strings"
)

type Config struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	DBName   string `yaml:"dbName"`
	SSLMode  bool   `yaml:"sslMode"`
	Password string `yaml:"password"`
}

func NewGorm(cfg *Config) (*gorm.DB, error) {

	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.DBName,
		cfg.Password,
	)

	err := createDB(cfg)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(gorm_postgres.Open(dataSourceName), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return db, nil
}

//Ref: https://dev.to/rafaelgfirmino/pagination-using-gorm-scopes-3k5f

func Paginate[T any](ctx context.Context, listQuery *utils.ListQuery, db *gorm.DB) (*utils.ListResult[T], error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "gorm.Paginate")

	var items []*T
	var totalRows int64
	db.Model(items).Count(&totalRows)

	// generate where query
	query := db.Offset(listQuery.GetOffset()).Limit(listQuery.GetLimit()).Order(listQuery.GetOrderBy())

	if listQuery.Filters != nil {
		for _, filter := range listQuery.Filters {
			column := filter.Field
			action := filter.Comparison
			value := filter.Value

			switch action {
			case "equals":
				whereQuery := fmt.Sprintf("%s = ?", column)
				query = query.Where(whereQuery, value)
				break
			case "contains":
				whereQuery := fmt.Sprintf("%s LIKE ?", column)
				query = query.Where(whereQuery, "%"+value+"%")
				break
			case "in":
				whereQuery := fmt.Sprintf("%s IN (?)", column)
				queryArray := strings.Split(value, ",")
				query = query.Where(whereQuery, queryArray)
				break

			}
		}
	}

	if err := query.Find(&items).Error; err != nil {
		tracing.TraceErr(span, err)
		return nil, errors.Wrap(err, "error in finding products.")
	}

	return utils.NewListResult(items, listQuery.GetSize(), listQuery.GetPage(), totalRows), nil
}

func createDB(cfg *Config) error {

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
