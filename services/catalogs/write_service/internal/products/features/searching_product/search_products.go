package searching_product

import "github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"

type SearchProducts struct {
	SearchText string `json:"searchText" validate:"required"`
	*utils.ListQuery
}
