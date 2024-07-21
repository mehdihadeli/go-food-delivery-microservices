package dtos

import (
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/utils"
)

type SearchProductsRequestDto struct {
	SearchText       string `query:"search" json:"search"`
	*utils.ListQuery `                      json:"listQuery"`
}
