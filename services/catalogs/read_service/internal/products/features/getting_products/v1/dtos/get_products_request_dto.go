package dtos

import "github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"

type GetProductsRequestDto struct {
	*utils.ListQuery
}
