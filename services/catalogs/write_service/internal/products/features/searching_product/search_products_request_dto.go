package searching_product

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
)

type SearchProductsRequestDto struct {
	SearchText       string `query:"search" json:"search" validate:"required"`
	*utils.ListQuery  `json:"listQuery"`
}
