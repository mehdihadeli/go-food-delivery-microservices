package endpoints

import (
	"net/http"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/web/route"
	customErrors "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/httperrors/customerrors"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/utils"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogreadservice/internal/products/contracts/params"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogreadservice/internal/products/features/searching_products/v1/dtos"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogreadservice/internal/products/features/searching_products/v1/queries"

	"emperror.dev/errors"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
)

type searchProductsEndpoint struct {
	params.ProductRouteParams
}

func NewSearchProductsEndpoint(
	params params.ProductRouteParams,
) route.Endpoint {
	return &searchProductsEndpoint{
		ProductRouteParams: params,
	}
}

func (ep *searchProductsEndpoint) MapEndpoint() {
	ep.ProductsGroup.GET("/search", ep.handler())
}

// SearchProducts
// @Tags Products
// @Summary Search products
// @Description Search products
// @Accept json
// @Produce json
// @Param searchProductsRequestDto query dtos.SearchProductsRequestDto false "SearchProductsRequestDto"
// @Success 200 {object} dtos.SearchProductsResponseDto
// @Router /api/v1/products/search [get]
func (ep *searchProductsEndpoint) handler() echo.HandlerFunc {
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

		request := &dtos.SearchProductsRequestDto{ListQuery: listQuery}

		if err := c.Bind(request); err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(
				err,
				"error in the binding request",
			)

			return badRequestErr
		}

		query := &queries.SearchProducts{
			SearchText: request.SearchText,
			ListQuery:  request.ListQuery,
		}

		if err := query.Validate(); err != nil {
			validationErr := customErrors.NewValidationErrorWrap(
				err,
				"query validation failed",
			)

			return validationErr
		}

		queryResult, err := mediatr.Send[*queries.SearchProducts, *dtos.SearchProductsResponseDto](
			ctx,
			query,
		)
		if err != nil {
			return errors.WithMessage(
				err,
				"error in sending SearchProducts",
			)
		}

		return c.JSON(http.StatusOK, queryResult)
	}
}
