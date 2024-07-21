package v1

import (
	"net/http"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/web/route"
	customErrors "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/httperrors/customerrors"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/dtos/v1/fxparams"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/features/creatingproduct/v1/dtos"

	"emperror.dev/errors"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
)

type createProductEndpoint struct {
	fxparams.ProductRouteParams
}

func NewCreteProductEndpoint(
	params fxparams.ProductRouteParams,
) route.Endpoint {
	return &createProductEndpoint{ProductRouteParams: params}
}

func (ep *createProductEndpoint) MapEndpoint() {
	ep.ProductsGroup.POST("", ep.handler())
}

// CreateProduct
// @Tags Products
// @Summary Create product
// @Description Create new product item
// @Accept json
// @Produce json
// @Param CreateProductRequestDto body dtos.CreateProductRequestDto true "Product data"
// @Success 201 {object} dtos.CreateProductResponseDto
// @Router /api/v1/products [post]
func (ep *createProductEndpoint) handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		request := &dtos.CreateProductRequestDto{}
		if err := c.Bind(request); err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(
				err,
				"error in the binding request",
			)

			return badRequestErr
		}

		command, err := NewCreateProductWithValidation(
			request.Name,
			request.Description,
			request.Price,
		)
		if err != nil {
			return err
		}

		result, err := mediatr.Send[*CreateProduct, *dtos.CreateProductResponseDto](
			ctx,
			command,
		)
		if err != nil {
			return errors.WithMessage(
				err,
				"error in sending CreateProduct",
			)
		}

		return c.JSON(http.StatusCreated, result)
	}
}
