package v1

import "github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"

// Ref: https://golangbot.com/inheritance/

type GetProductsQuery struct {
	*utils.ListQuery
}

func NewGetProductsQuery(query *utils.ListQuery) *GetProductsQuery {
	return &GetProductsQuery{ListQuery: query}
}
