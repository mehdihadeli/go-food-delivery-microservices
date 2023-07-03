package endpoints

import (
	"fmt"
	"net/http"

	"emperror.dev/errors"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/contracts/params"

	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/web/route"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/updating_product/v1/commands"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/updating_product/v1/dtos"
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
				"[updateProductEndpoint_handler.Bind] error in the binding request",
			)
			ep.Logger.Errorf(
				fmt.Sprintf("[updateProductEndpoint_handler.Bind] err: %v", badRequestErr),
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
				"[updateProductEndpoint_handler.StructCtx] command validation failed",
			)
			ep.Logger.Errorf(
				fmt.Sprintf("[updateProductEndpoint_handler.StructCtx] err: {%v}", validationErr),
			)
			return validationErr
		}

		_, err = mediatr.Send[*commands.UpdateProduct, *mediatr.Unit](ctx, command)
		if err != nil {
			err = errors.WithMessage(
				err,
				"[updateProductEndpoint_handler.Send] error in sending UpdateProduct",
			)
			ep.Logger.Errorw(
				fmt.Sprintf(
					"[updateProductEndpoint_handler.Send] id: {%s}, err: {%v}",
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
