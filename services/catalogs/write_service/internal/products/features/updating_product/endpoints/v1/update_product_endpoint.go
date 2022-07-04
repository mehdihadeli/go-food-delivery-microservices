package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/updating_product"
	shared_configurations "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/configurations"
	"net/http"
)

type updateProductEndpoint struct {
	mediator          *mediatr.Mediator
	productRepository contracts.ProductRepository
	productsGroup     *echo.Group
	infrastructure    *shared_configurations.Infrastructure
}

func NewUpdateProductEndpoint(infra *shared_configurations.Infrastructure, mediator *mediatr.Mediator, productsGroup *echo.Group, productRepository contracts.ProductRepository) *updateProductEndpoint {
	return &updateProductEndpoint{mediator: mediator, productRepository: productRepository, productsGroup: productsGroup, infrastructure: infra}
}

func (ep *updateProductEndpoint) MapRoute() {
	ep.productsGroup.PUT("/:id", ep.updateProduct())
}

// UpdateProduct
// @Tags Products
// @Summary Update product
// @Description Update existing product
// @Accept json
// @Produce json
// @Param UpdateProductRequestDto body updating_product.UpdateProductRequestDto true "Product data"
// @Param id path string true "Product ID"
// @Success 204
// @Router /products/{id} [put]
func (ep *updateProductEndpoint) updateProduct() echo.HandlerFunc {
	return func(c echo.Context) error {

		ep.infrastructure.Metrics.UpdateProductHttpRequests.Inc()
		ctx, span := tracing.StartHttpServerTracerSpan(c, "updateProductEndpoint.updateProduct")
		defer span.Finish()

		request := &updating_product.UpdateProductRequestDto{}
		if err := c.Bind(request); err != nil {
			ep.infrastructure.Log.WarnMsg("Bind", err)
			tracing.TraceErr(span, err)
			return err
		}

		command := updating_product.NewUpdateProduct(request.ProductID, request.Name, request.Description, request.Price)

		if err := ep.infrastructure.Validator.StructCtx(ctx, command); err != nil {
			ep.infrastructure.Log.WarnMsg("validate", err)
			tracing.TraceErr(span, err)
			return err
		}

		_, err := ep.mediator.Send(ctx, command)

		if err != nil {
			ep.infrastructure.Log.WarnMsg("UpdateProduct", err)
			tracing.TraceErr(span, err)
			return err
		}

		ep.infrastructure.Log.Infof("(product updated) id: {%s}", request.ProductID)

		return c.NoContent(http.StatusNoContent)
	}
}
