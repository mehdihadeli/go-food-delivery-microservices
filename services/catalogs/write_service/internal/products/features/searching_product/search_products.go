package searching_product

import "github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"

type SearchProducts struct {
	SearchText string `validate:"required"`
	*utils.ListQuery
}
