package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/deleting_product"
	shared_configurations "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/configurations"
	"net/http"
)

type deleteProductEndpoint struct {
	mediator          *mediatr.Mediator
	productRepository contracts.ProductRepository
	productsGroup     *echo.Group
	infrastructure    *shared_configurations.Infrastructure
}

func NewDeleteProductEndpoint(infra *shared_configurations.Infrastructure, mediator *mediatr.Mediator, productsGroup *echo.Group, productRepository contracts.ProductRepository) *deleteProductEndpoint {
	return &deleteProductEndpoint{mediator: mediator, productRepository: productRepository, productsGroup: productsGroup, infrastructure: infra}
}

func (ep *deleteProductEndpoint) MapRoute() {
	ep.productsGroup.DELETE("/:id", ep.deleteProduct())
}

// DeleteProduct
// @Tags Products
// @Summary Delete product
// @Description Delete existing product
// @Accept json
// @Produce json
// @Success 204
// @Param id path string true "Product ID"
// @Router /products/{id} [delete]
func (ep *deleteProductEndpoint) deleteProduct() echo.HandlerFunc {
	return func(c echo.Context) error {

		ep.infrastructure.Metrics.DeleteProductHttpRequests.Inc()
		ctx, span := tracing.StartHttpServerTracerSpan(c, "deleteProductEndpoint.deleteProduct")
		defer span.Finish()

		request := &deleting_product.DeleteProductRequestDto{}
		if err := c.Bind(request); err != nil {
			ep.infrastructure.Log.WarnMsg("Bind", err)
			tracing.TraceErr(span, err)
			return err
		}

		command := deleting_product.NewDeleteProduct(request.ProductID)

		if err := ep.infrastructure.Validator.StructCtx(ctx, command); err != nil {
			ep.infrastructure.Log.WarnMsg("validate", err)
			tracing.TraceErr(span, err)
			return err
		}

		_, err := ep.mediator.Send(ctx, command)

		if err != nil {
			ep.infrastructure.Log.WarnMsg("DeleteProduct", err)
			tracing.TraceErr(span, err)
			return err
		}

		return c.NoContent(http.StatusNoContent)
	}
}
