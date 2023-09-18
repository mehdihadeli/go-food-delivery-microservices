package gormPostgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils"

	"emperror.dev/errors"
	"gorm.io/gorm"
)

// Ref: https://dev.to/rafaelgfirmino/pagination-using-gorm-scopes-3k5f

func Paginate[T any](
	ctx context.Context,
	listQuery *utils.ListQuery,
	db *gorm.DB,
) (*utils.ListResult[T], error) {
	var items []T
	var totalRows int64
	db.Model(items).WithContext(ctx).Count(&totalRows)

	// generate where query
	query := db.WithContext(ctx).
		Offset(listQuery.GetOffset()).
		Limit(listQuery.GetLimit()).
		Order(listQuery.GetOrderBy())

	if listQuery.Filters != nil {
		for _, filter := range listQuery.Filters {
			column := filter.Field
			action := filter.Comparison
			value := filter.Value

			switch action {
			case "equals":
				whereQuery := fmt.Sprintf("%s = ?", column)
				query = query.WithContext(ctx).Where(whereQuery, value)
				break
			case "contains":
				whereQuery := fmt.Sprintf("%s LIKE ?", column)
				query = query.WithContext(ctx).Where(whereQuery, "%"+value+"%")
				break
			case "in":
				whereQuery := fmt.Sprintf("%s IN (?)", column)
				queryArray := strings.Split(value, ",")
				query = query.WithContext(ctx).Where(whereQuery, queryArray)
				break

			}
		}
	}

	if err := query.Find(&items).Error; err != nil {
		return nil, errors.WrapIf(err, "error in finding products.")
	}

	return utils.NewListResult[T](items, listQuery.GetSize(), listQuery.GetPage(), totalRows), nil
}
