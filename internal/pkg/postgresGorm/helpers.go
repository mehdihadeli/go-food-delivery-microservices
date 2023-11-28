package postgresGorm

import (
	"context"
	"fmt"
	"strings"

	typeMapper "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/reflection/typemapper"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils"

	"emperror.dev/errors"
	"gorm.io/gorm"
)

// Ref: https://dev.to/rafaelgfirmino/pagination-using-gorm-scopes-3k5f

func Paginate[TDataModel any, TEntity any](
	ctx context.Context,
	listQuery *utils.ListQuery,
	db *gorm.DB,
) (*utils.ListResult[TEntity], error) {
	var (
		items     []TEntity
		totalRows int64
	)

	dataModel := typeMapper.GenericInstanceByT[TDataModel]()
	// https://gorm.io/docs/advanced_query.html
	db.WithContext(ctx).Model(dataModel).Count(&totalRows)

	// generate where query
	query := db.WithContext(ctx).
		Model(dataModel).
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
				query = query.Where(whereQuery, value)
			case "contains":
				whereQuery := fmt.Sprintf("%s LIKE ?", column)
				query = query.Where(whereQuery, "%"+value+"%")
			case "in":
				whereQuery := fmt.Sprintf("%s IN (?)", column)
				queryArray := strings.Split(value, ",")
				query = query.Where(whereQuery, queryArray)
			}
		}
	}

	// https://gorm.io/docs/advanced_query.html#Smart-Select-Fields
	if err := query.Find(&items).Error; err != nil {
		return nil, errors.WrapIf(err, "error in finding products.")
	}

	return utils.NewListResult[TEntity](
		items,
		listQuery.GetSize(),
		listQuery.GetPage(),
		totalRows,
	), nil
}
