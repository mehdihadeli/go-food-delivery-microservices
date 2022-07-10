package dtos

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/dto"
)

type SearchProductsResponseDto struct {
	Products *utils.ListResult[*dto.ProductDto]
}
