package dtos

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	dtoV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/dto/v1"
)

type SearchProductsResponseDto struct {
	Products *utils.ListResult[*dtoV1.ProductDto]
}
