package gorm_postgres

import (
	"fmt"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"math"
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

	db, err := gorm.Open(postgres.Open(dataSourceName), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return db, nil
}

//Ref: https://dev.to/rafaelgfirmino/pagination-using-gorm-scopes-3k5f

func Paginate[T any](pagination *utils.ListQuery, db *gorm.DB) (*utils.ListResult[T], error) {

	var items []*T
	var totalRows int64
	db.Model(items).Count(&totalRows)

	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.GetSize())))

	// generate where query
	query := db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination.GetOrderBy())

	if pagination.Filters != nil {
		for _, filter := range pagination.Filters {
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

	if result := query.Find(&items); result.Error != nil {
		return nil, errors.Wrap(result.Error, "error in finding products.")
	}

	return utils.NewListResult(items, pagination.GetSize(), pagination.GetPage(), totalRows, totalPages), nil
}
