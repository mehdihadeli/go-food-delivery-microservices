package dtos

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils"
)

type SearchProductsRequestDto struct {
	SearchText       string `query:"search" json:"search"`
	*utils.ListQuery `                      json:"listQuery"`
}
