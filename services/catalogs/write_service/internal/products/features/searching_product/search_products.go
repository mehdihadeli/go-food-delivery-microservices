package searching_product

import "github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"

type SearchProducts struct {
	SearchText string
	*utils.ListQuery
}
