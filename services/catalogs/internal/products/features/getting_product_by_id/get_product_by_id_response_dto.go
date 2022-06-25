package getting_product_by_id

import "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/dto"

type GetProductByIdResponseDto struct {
	Product *dto.ProductDto
}
