package v1

import (
	"emperror.dev/errors"
	"fmt"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/deleting_product/commands/v1"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/delivery"
	deletingProduct "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/deleting_product"
)

type deleteProductEndpoint struct {
	*delivery.ProductEndpointBase
}

func NewDeleteProductEndpoint(productEndpointBase *delivery.ProductEndpointBase) *deleteProductEndpoint {
	return &deleteProductEndpoint{productEndpointBase}
}

func (ep *deleteProductEndpoint) MapRoute() {
	ep.ProductsGroup.DELETE("/:id", ep.handler())
}

// DeleteProductCommand
// @Tags Products
// @Summary Delete product
// @Description Delete existing product
// @Accept json
// @Produce json
// @Success 204
// @Param id path string true "Product ID"
// @Router /api/v1/products/{id} [delete]
func (ep *deleteProductEndpoint) handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		ep.Metrics.DeleteProductHttpRequests.Inc()
		ctx, span := tracing.StartHttpServerTracerSpan(c, "deleteProductEndpoint.handler")
		defer span.Finish()

		request := &deletingProduct.DeleteProductRequestDto{}
		if err := c.Bind(request); err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(err, "[deleteProductEndpoint_handler.Bind] error in the binding request")
			ep.Log.Errorf(fmt.Sprintf("[deleteProductEndpoint_handler.Bind] err: %v", tracing.TraceWithErr(span, badRequestErr)))
			return badRequestErr
		}

		command := v1.NewDeleteProductCommand(request.ProductID)
		if err := ep.Validator.StructCtx(ctx, command); err != nil {
			validationErr := customErrors.NewValidationErrorWrap(err, "[deleteProductEndpoint_handler.StructCtx] command validation failed")
			ep.Log.Errorf(fmt.Sprintf("[deleteProductEndpoint_handler.StructCtx] err: {%v}", tracing.TraceWithErr(span, validationErr)))
			return validationErr
		}

		_, err := mediatr.Send[*v1.DeleteProductCommand, *mediatr.Unit](ctx, command)

		if err != nil {
			err = errors.WithMessage(err, "[deleteProductEndpoint_handler.Send] error in sending DeleteProductCommand")
			ep.Log.Errorw(fmt.Sprintf("[deleteProductEndpoint_handler.Send] id: {%s}, err: {%v}", command.ProductID, tracing.TraceWithErr(span, err)), logger.Fields{"ProductId": command.ProductID})
			return err
		}

		return c.NoContent(http.StatusNoContent)
	}
}
