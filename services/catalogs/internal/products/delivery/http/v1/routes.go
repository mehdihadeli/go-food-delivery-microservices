package v1

//
//import (
//	"github.com/labstack/echo/v4"
//	"net/http"
//)
//
//func (h *productsController) MapRoutes() {
//
//	// ref: https://dev.to/krishnakummar/api-versioning-in-golang-echo-5eh5
//	// https://medium.com/onexlab/go-lang-echo-framework-header-based-api-versioning-9a6701f6d38
//	v1 := h.e.Group("/api/v1")
//	groupV1Routes(v1, h)
//}
//
//func groupV1Routes(group *echo.Group, h *productsController) {
//
//	products := group.Group("/" + h.cfg.Http.ProductsPath)
//	products.GET("", h.GetAllProducts())
//	products.POST("", h.CreateProduct())
//	products.GET("/:id", h.GetProductByID())
//	products.PUT("/:id", h.UpdateProduct())
//	products.DELETE("/:id", h.DeleteProduct())
//	products.Any("/health", func(c echo.Context) error {
//		return c.JSON(http.StatusOK, "OK")
//	})
//}
