package endpoints

import (
	"net/http"

	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/web/route"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/contracts/params"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/updating_product/v1/commands"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/updating_product/v1/dtos"

	"emperror.dev/errors"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
)

type updateProductEndpoint struct {
	params.ProductRouteParams
}

func NewUpdateProductEndpoint(
	params params.ProductRouteParams,
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
		ep.CatalogsMetrics.UpdateProductHttpRequests.Add(ctx, 1)

		request := &dtos.UpdateProductRequestDto{}
		if err := c.Bind(request); err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(
				err,
				"error in the binding request",
			)

			return badRequestErr
		}

		command, err := commands.NewUpdateProduct(
			request.ProductID,
			request.Name,
			request.Description,
			request.Price,
		)
		if err != nil {
			validationErr := customErrors.NewValidationErrorWrap(
				err,
				"command validation failed",
			)

			return validationErr
		}

		_, err = mediatr.Send[*commands.UpdateProduct, *mediatr.Unit](ctx, command)
		if err != nil {
			return errors.WithMessage(
				err,
				"error in sending UpdateProduct",
			)
		}

		return c.NoContent(http.StatusNoContent)
	}
}
