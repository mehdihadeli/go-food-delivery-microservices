package dtos

import "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/products/dto"

type GetProductByIdResponseDto struct {
	Product *dto.ProductDto `json:"product"`
}
