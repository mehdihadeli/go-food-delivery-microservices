package v1

import (
	"fmt"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/updating_product/commands/v1"
	"github.com/pkg/errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/delivery"
	updatingProduct "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/updating_product"
)

type updateProductEndpoint struct {
	*delivery.ProductEndpointBase
}

func NewUpdateProductEndpoint(productEndpointBase *delivery.ProductEndpointBase) *updateProductEndpoint {
	return &updateProductEndpoint{productEndpointBase}
}

func (ep *updateProductEndpoint) MapRoute() {
	ep.ProductsGroup.PUT("/:id", ep.handler())
}

// UpdateProductCommand
// @Tags Products
// @Summary Update product
// @Description Update existing product
// @Accept json
// @Produce json
// @Param UpdateProductRequestDto body updatingProduct.UpdateProductRequestDto true "Product data"
// @Param id path string true "Product ID"
// @Success 204
// @Router /api/v1/products/{id} [put]
func (ep *updateProductEndpoint) handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		ep.Metrics.UpdateProductHttpRequests.Inc()
		ctx, span := tracing.StartHttpServerTracerSpan(c, "updateProductEndpoint.updateProduct")
		defer span.Finish()

		request := &updatingProduct.UpdateProductRequestDto{}
		if err := c.Bind(request); err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(err, "[updateProductEndpoint_handler.Bind] error in the binding request")
			ep.Log.Errorf(fmt.Sprintf("[updateProductEndpoint_handler.Bind] err: %v", tracing.TraceWithErr(span, badRequestErr)))
			return badRequestErr
		}

		command := v1.NewUpdateProductCommand(request.ProductID, request.Name, request.Description, request.Price)

		if err := ep.Validator.StructCtx(ctx, command); err != nil {
			validationErr := customErrors.NewValidationErrorWrap(err, "[updateProductEndpoint_handler.StructCtx] command validation failed")
			ep.Log.Errorf(fmt.Sprintf("[updateProductEndpoint_handler.StructCtx] err: {%v}", tracing.TraceWithErr(span, validationErr)))
			return validationErr
		}

		_, err := mediatr.Send[*v1.UpdateProductCommand, *mediatr.Unit](ctx, command)

		if err != nil {
			err = errors.WithMessage(err, "[updateProductEndpoint_handler.Send] error in sending UpdateProductCommand")
			ep.Log.Errorw(fmt.Sprintf("[updateProductEndpoint_handler.Send] id: {%s}, err: {%v}", command.ProductID, tracing.TraceWithErr(span, err)), logger.Fields{"ProductId": command.ProductID})
			return err
		}

		return c.NoContent(http.StatusNoContent)
	}
}
