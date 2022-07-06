package dtos

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/dto"
)

type GetProductsResponseDto struct {
	Products *utils.ListResult[dto.ProductDto]
}
