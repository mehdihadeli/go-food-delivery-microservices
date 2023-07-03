package dtos

// https://echo.labstack.com/guide/binding/
// https://echo.labstack.com/guide/request/
// https://github.com/go-playground/validator

// CreateProductRequestDto validation will handle in command level
type CreateProductRequestDto struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}
