package searching_products

import "github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"

type SearchProducts struct {
	SearchText string `validate:"required"`
	*utils.ListQuery
}
