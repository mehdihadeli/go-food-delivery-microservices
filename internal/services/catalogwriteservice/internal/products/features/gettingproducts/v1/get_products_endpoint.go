package v1

import (
	"net/http"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/web/route"
	customErrors "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/httperrors/customerrors"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/utils"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/dtos/v1/fxparams"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/features/gettingproducts/v1/dtos"

	"emperror.dev/errors"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
)

type getProductsEndpoint struct {
	fxparams.ProductRouteParams
}

func NewGetProductsEndpoint(
	params fxparams.ProductRouteParams,
) route.Endpoint {
	return &getProductsEndpoint{ProductRouteParams: params}
}

func (ep *getProductsEndpoint) MapEndpoint() {
	ep.ProductsGroup.GET("", ep.handler())
}

// GetAllProducts
// @Tags Products
// @Summary Get all product
// @Description Get all products
// @Accept json
// @Produce json
// @Param getProductsRequestDto query dtos.GetProductsRequestDto false "GetProductsRequestDto"
// @Success 200 {object} dtos.GetProductsResponseDto
// @Router /api/v1/products [get]
func (ep *getProductsEndpoint) handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		listQuery, err := utils.GetListQueryFromCtx(c)
		if err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(
				err,
				"error in getting data from query string",
			)

			return badRequestErr
		}

		request := &dtos.GetProductsRequestDto{ListQuery: listQuery}
		if err := c.Bind(request); err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(
				err,
				"error in the binding request",
			)

			return badRequestErr
		}

		query, err := NewGetProducts(request.ListQuery)
		if err != nil {
			return err
		}

		queryResult, err := mediatr.Send[*GetProducts, *dtos.GetProductsResponseDto](
			ctx,
			query,
		)
		if err != nil {
			return errors.WithMessage(
				err,
				"error in sending GetProducts",
			)
		}

		return c.JSON(http.StatusOK, queryResult)
	}
}
