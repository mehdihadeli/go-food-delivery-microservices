package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/constants"
	httpErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts/repositories"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/deleting_product"
	shared_configurations "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/configurations"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

type deleteProductEndpoint struct {
	mediator          *mediatr.Mediator
	productRepository repositories.ProductRepository
	productsGroup     *echo.Group
	infrastructure    *shared_configurations.Infrastructure
}

func NewDeleteProductEndpoint(infra *shared_configurations.Infrastructure, mediator *mediatr.Mediator, productsGroup *echo.Group, productRepository repositories.ProductRepository) *deleteProductEndpoint {
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

		productUUID, err := uuid.FromString(c.Param(constants.ID))
		if err != nil {
			ep.infrastructure.Log.WarnMsg("uuid.FromString", err)
			ep.infrastructure.TraceErr(span, err)
			return httpErrors.ErrorResponse(err, ep.infrastructure.Cfg.Http.DebugErrorsResponse)
		}

		command := deleting_product.NewDeleteProduct(productUUID)
		_, err = ep.mediator.Send(ctx, command)

		if err != nil {
			ep.infrastructure.Log.WarnMsg("DeleteProduct", err)
			ep.infrastructure.Metrics.ErrorHttpRequests.Inc()
			return httpErrors.ErrorResponse(err, ep.infrastructure.Cfg.Http.DebugErrorsResponse)
		}

		ep.infrastructure.Metrics.SuccessHttpRequests.Inc()
		return c.NoContent(http.StatusNoContent)
	}
}
