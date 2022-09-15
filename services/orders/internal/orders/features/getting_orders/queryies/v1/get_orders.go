package v1

import "github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"

// Ref: https://golangbot.com/inheritance/

type GetOrders struct {
	*utils.ListQuery
}

func NewGetOrders(query *utils.ListQuery) *GetOrders {
	return &GetOrders{ListQuery: query}
}
