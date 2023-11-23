package v1

import (
	"net/http"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/web/route"
	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/dtos/v1/fxparams"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/deleting_product/v1/dtos"

	"emperror.dev/errors"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
)

type deleteProductEndpoint struct {
	fxparams.ProductRouteParams
}

func NewDeleteProductEndpoint(
	params fxparams.ProductRouteParams,
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

		request := &dtos.DeleteProductRequestDto{}
		if err := c.Bind(request); err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(
				err,
				"error in the binding request",
			)

			return badRequestErr
		}

		command, err := NewDeleteProduct(request.ProductID)
		if err != nil {
			validationErr := customErrors.NewValidationErrorWrap(
				err,
				"command validation failed",
			)

			return validationErr
		}

		_, err = mediatr.Send[*DeleteProduct, *mediatr.Unit](
			ctx,
			command,
		)

		if err != nil {
			return errors.WithMessage(
				err,
				"error in sending DeleteProduct",
			)
		}

		return c.NoContent(http.StatusNoContent)
	}
}
