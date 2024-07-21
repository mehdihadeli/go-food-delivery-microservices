package dtos

import "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/utils"

type GetProductsRequestDto struct {
	*utils.ListQuery
}
