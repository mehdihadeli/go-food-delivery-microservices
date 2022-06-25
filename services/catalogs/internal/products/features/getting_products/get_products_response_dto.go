package getting_products

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/dto"
)

type GetProductsResponseDto struct {
	Products []*dto.ProductDto
}
