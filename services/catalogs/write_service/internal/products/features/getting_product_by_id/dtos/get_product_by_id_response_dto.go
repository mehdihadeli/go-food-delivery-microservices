package dtos

import "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/dto"

type GetProductByIdResponseDto struct {
	Product *dto.ProductDto `json:"product"`
}
