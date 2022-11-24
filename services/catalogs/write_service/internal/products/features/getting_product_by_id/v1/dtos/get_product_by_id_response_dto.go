package dtos

import dtoV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/dto/v1"

// https://echo.labstack.com/guide/response/
type GetProductByIdResponseDto struct {
	Product *dtoV1.ProductDto `json:"product"`
}
