package endpoints

import (
	"fmt"
	"net/http"

	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/web/route"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/contracts/params"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/deleting_product/v1/commands"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/deleting_product/v1/dtos"

	"emperror.dev/errors"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
)

type deleteProductEndpoint struct {
	params.ProductRouteParams
}

func NewDeleteProductEndpoint(
	params params.ProductRouteParams,
) route.Endpoint {
	return &deleteProductEndpoint{ProductRouteParams: params}
}

func (ep *deleteProductEndpoint) MapEndpoint() {
	ep.ProductsGroup.DELETE("/:id", ep.handler())
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
func (ep *deleteProductEndpoint) handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		ep.CatalogsMetrics.DeleteProductHttpRequests.Add(ctx, 1)

		request := &dtos.DeleteProductRequestDto{}
		if err := c.Bind(request); err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(
				err,
				"[deleteProductEndpoint_handler.Bind] error in the binding request",
			)
			ep.Logger.Errorf(
				fmt.Sprintf("[deleteProductEndpoint_handler.Bind] err: %v", badRequestErr),
			)
			return badRequestErr
		}

		command, err := commands.NewDeleteProduct(request.ProductID)
		if err != nil {
			validationErr := customErrors.NewValidationErrorWrap(
				err,
				"[deleteProductEndpoint_handler.StructCtx] command validation failed",
			)
			ep.Logger.Errorf(
				fmt.Sprintf("[deleteProductEndpoint_handler.StructCtx] err: {%v}", validationErr),
			)
			return validationErr
		}

		_, err = mediatr.Send[*commands.DeleteProduct, *mediatr.Unit](ctx, command)

		if err != nil {
			err = errors.WithMessage(
				err,
				"[deleteProductEndpoint_handler.Send] error in sending DeleteProduct",
			)
			ep.Logger.Errorw(
				fmt.Sprintf(
					"[deleteProductEndpoint_handler.Send] id: {%s}, err: {%v}",
					command.ProductID,
					err,
				),
				logger.Fields{"ProductId": command.ProductID},
			)
			return err
		}

		return c.NoContent(http.StatusNoContent)
	}
}
