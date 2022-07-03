package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts/repositories"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/updating_product"
	shared_configurations "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/configurations"
	"net/http"
)

type updateProductEndpoint struct {
	mediator          *mediatr.Mediator
	productRepository repositories.ProductRepository
	productsGroup     *echo.Group
	infrastructure    *shared_configurations.Infrastructure
}

func NewUpdateProductEndpoint(infra *shared_configurations.Infrastructure, mediator *mediatr.Mediator, productsGroup *echo.Group, productRepository repositories.ProductRepository) *updateProductEndpoint {
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
			ep.infrastructure.TraceErr(span, err)
			return err
		}

		if err := ep.infrastructure.Validator.StructCtx(ctx, request); err != nil {
			ep.infrastructure.Log.WarnMsg("validate", err)
			ep.infrastructure.TraceErr(span, err)
			return err
		}

		command := updating_product.NewUpdateProduct(request.ProductID, request.Name, request.Description, request.Price)

		_, err := ep.mediator.Send(ctx, command)

		if err != nil {
			ep.infrastructure.Log.WarnMsg("UpdateProduct", err)
			ep.infrastructure.Metrics.ErrorHttpRequests.Inc()
			return err
		}

		ep.infrastructure.Log.Infof("(product updated) id: {%s}", request.ProductID)
		ep.infrastructure.Metrics.SuccessHttpRequests.Inc()

		return c.NoContent(http.StatusNoContent)
	}
}
