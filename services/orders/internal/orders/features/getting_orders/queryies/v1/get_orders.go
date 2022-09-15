package v1

import "github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"

// Ref: https://golangbot.com/inheritance/

type GetOrdersQuery struct {
	*utils.ListQuery
}

func NewGetOrdersQuery(query *utils.ListQuery) *GetOrdersQuery {
	return &GetOrdersQuery{ListQuery: query}
}
