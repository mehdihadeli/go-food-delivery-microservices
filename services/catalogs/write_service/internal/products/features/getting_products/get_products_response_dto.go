package getting_products

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/dto"
)

type GetProductsResponseDto struct {
	Products []*dto.ProductDto
}
