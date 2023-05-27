package dtos

import "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/products/dto"

type GetProductByIdResponseDto struct {
	Product *dto.ProductDto `json:"product"`
}
