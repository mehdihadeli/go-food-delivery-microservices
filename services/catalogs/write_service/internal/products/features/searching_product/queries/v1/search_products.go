package v1

import "github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"

type SearchProducts struct {
	SearchText string `validate:"required"`
	*utils.ListQuery
}

func NewSearchProducts(searchText string, query *utils.ListQuery) *SearchProducts {
	return &SearchProducts{
		SearchText: searchText,
		ListQuery:  query,
	}
}
