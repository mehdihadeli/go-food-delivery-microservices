package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/constants"
	httpErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/contracts/repositories"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/features/updating_product"
	shared_configurations "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/shared/configurations"
	uuid "github.com/satori/go.uuid"
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

		productUUID, err := uuid.FromString(c.Param(constants.ID))
		if err != nil {
			ep.infrastructure.Log.WarnMsg("uuid.FromString", err)
			ep.infrastructure.TraceErr(span, err)
			return httpErrors.ErrorCtxResponse(c, err, ep.infrastructure.Cfg.Http.DebugErrorsResponse)
		}

		request := &updating_product.UpdateProductRequestDto{ProductID: productUUID}
		if err := c.Bind(request); err != nil {
			ep.infrastructure.Log.WarnMsg("Bind", err)
			ep.infrastructure.TraceErr(span, err)
			return httpErrors.ErrorCtxResponse(c, err, ep.infrastructure.Cfg.Http.DebugErrorsResponse)
		}

		if err := ep.infrastructure.Validator.StructCtx(ctx, request); err != nil {
			ep.infrastructure.Log.WarnMsg("validate", err)
			ep.infrastructure.TraceErr(span, err)
			return httpErrors.ErrorCtxResponse(c, err, ep.infrastructure.Cfg.Http.DebugErrorsResponse)
		}

		command := updating_product.NewUpdateProduct(productUUID, request.Name, request.Description, request.Price)

		_, err = ep.mediator.Send(ctx, command)

		if err != nil {
			ep.infrastructure.Log.WarnMsg("UpdateProduct", err)
			ep.infrastructure.Metrics.ErrorHttpRequests.Inc()
			return httpErrors.ErrorCtxResponse(c, err, ep.infrastructure.Cfg.Http.DebugErrorsResponse)
		}

		ep.infrastructure.Log.Infof("(product updated) id: {%s}", productUUID.String())
		ep.infrastructure.Metrics.SuccessHttpRequests.Inc()

		return c.NoContent(http.StatusNoContent)
	}
}
