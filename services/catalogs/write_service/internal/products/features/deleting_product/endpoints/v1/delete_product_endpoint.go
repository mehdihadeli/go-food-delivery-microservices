package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/delivery"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/deleting_product"
	"net/http"
)

type deleteProductEndpoint struct {
	*delivery.ProductEndpointBase
}

func NewDeleteProductEndpoint(productEndpointBase *delivery.ProductEndpointBase) *deleteProductEndpoint {
	return &deleteProductEndpoint{productEndpointBase}
}

func (ep *deleteProductEndpoint) MapRoute() {
	ep.ProductsGroup.DELETE("/:id", ep.deleteProduct())
}

// DeleteProduct
// @Tags Products
// @Summary Delete product
// @Description Delete existing product
// @Accept json
// @Produce json
// @Success 204
// @Param id path string true "Product ID"
// @Router /api/v1/products/{id} [delete]
func (ep *deleteProductEndpoint) deleteProduct() echo.HandlerFunc {
	return func(c echo.Context) error {

		ep.Metrics.DeleteProductHttpRequests.Inc()
		ctx, span := tracing.StartHttpServerTracerSpan(c, "deleteProductEndpoint.deleteProduct")
		defer span.Finish()

		request := &deleting_product.DeleteProductRequestDto{}
		if err := c.Bind(request); err != nil {
			ep.Log.WarnMsg("Bind", err)
			tracing.TraceErr(span, err)
			return err
		}

		command := deleting_product.NewDeleteProduct(request.ProductID)

		if err := ep.Validator.StructCtx(ctx, command); err != nil {
			ep.Log.WarnMsg("validate", err)
			tracing.TraceErr(span, err)
			return err
		}

		_, err := ep.ProductMediator.Send(ctx, command)

		if err != nil {
			ep.Log.WarnMsg("DeleteProduct", err)
			tracing.TraceErr(span, err)
			return err
		}

		return c.NoContent(http.StatusNoContent)
	}
}
