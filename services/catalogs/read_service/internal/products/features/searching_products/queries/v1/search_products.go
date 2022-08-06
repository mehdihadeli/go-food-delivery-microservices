package v1

import "github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"

type SearchProductsQuery struct {
	SearchText string `validate:"required"`
	*utils.ListQuery
}
