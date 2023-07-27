package dtos

import "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils"

type GetOrdersRequestDto struct {
	*utils.ListQuery
}
