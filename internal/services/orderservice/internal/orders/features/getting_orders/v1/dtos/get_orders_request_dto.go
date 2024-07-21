package dtos

import "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/utils"

type GetOrdersRequestDto struct {
	*utils.ListQuery
}
