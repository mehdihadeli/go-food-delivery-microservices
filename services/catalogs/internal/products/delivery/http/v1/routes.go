package v1

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *productsHandlers) MapRoutes() {
	h.group.POST("", h.CreateProduct())
	h.group.GET("/:id", h.GetProductByID())
	h.group.PUT("/:id", h.UpdateProduct())
	h.group.DELETE("/:id", h.DeleteProduct())
	h.group.Any("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})
}
