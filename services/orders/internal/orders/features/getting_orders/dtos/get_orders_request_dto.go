package dtos

import "github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"

type GetOrdersRequestDto struct {
	*utils.ListQuery
}
