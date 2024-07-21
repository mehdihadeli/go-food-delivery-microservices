package v1

import (
	"net/http"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/web/route"
	customErrors "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/httperrors/customerrors"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/dtos/v1/fxparams"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/features/updatingproduct/v1/dtos"

	"emperror.dev/errors"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
)

type updateProductEndpoint struct {
	fxparams.ProductRouteParams
}

func NewUpdateProductEndpoint(
	params fxparams.ProductRouteParams,
) route.Endpoint {
	return &updateProductEndpoint{ProductRouteParams: params}
}

func (ep *updateProductEndpoint) MapEndpoint() {
	ep.ProductsGroup.PUT("/:id", ep.handler())
}

// UpdateProduct
// @Tags Products
// @Summary Update product
// @Description Update existing product
// @Accept json
// @Produce json
// @Param UpdateProductRequestDto body dtos.UpdateProductRequestDto true "Product data"
// @Param id path string true "Product ID"
// @Success 204
// @Router /api/v1/products/{id} [put]
func (ep *updateProductEndpoint) handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		request := &dtos.UpdateProductRequestDto{}
		if err := c.Bind(request); err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(
				err,
				"error in the binding request",
			)

			return badRequestErr
		}

		command, err := NewUpdateProductWithValidation(
			request.ProductID,
			request.Name,
			request.Description,
			request.Price,
		)
		if err != nil {
			return err
		}

		_, err = mediatr.Send[*UpdateProduct, *mediatr.Unit](
			ctx,
			command,
		)
		if err != nil {
			return errors.WithMessage(
				err,
				"error in sending UpdateProduct",
			)
		}

		return c.NoContent(http.StatusNoContent)
	}
}
