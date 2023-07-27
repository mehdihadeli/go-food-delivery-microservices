package queries

import "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils"

// Ref: https://golangbot.com/inheritance/

type GetProducts struct {
	*utils.ListQuery
}

func NewGetProducts(query *utils.ListQuery) *GetProducts {
	return &GetProducts{ListQuery: query}
}
