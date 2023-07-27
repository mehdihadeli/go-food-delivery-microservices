package endpoints

import (
	"fmt"
	"net/http"

	"emperror.dev/errors"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/products/contracts/params"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/products/features/searching_products/v1/dtos"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/products/features/searching_products/v1/queries"

	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/web/route"
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
		ep.CatalogsMetrics.SearchProductHttpRequests.Add(ctx, 1)

		listQuery, err := utils.GetListQueryFromCtx(c)
		if err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(
				err,
				"[searchProductsEndpoint_handler.GetListQueryFromCtx] error in getting data from query string",
			)
			ep.Logger.Errorf(
				fmt.Sprintf(
					"[searchProductsEndpoint_handler.GetListQueryFromCtx] err: %v",
					badRequestErr,
				),
			)
			return err
		}

		request := &dtos.SearchProductsRequestDto{ListQuery: listQuery}

		if err := c.Bind(request); err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(
				err,
				"[searchProductsEndpoint_handler.Bind] error in the binding request",
			)
			ep.Logger.Errorf(
				fmt.Sprintf("[searchProductsEndpoint_handler.Bind] err: %v", badRequestErr),
			)
			return badRequestErr
		}

		query := &queries.SearchProducts{
			SearchText: request.SearchText,
			ListQuery:  request.ListQuery,
		}

		if err := ep.Validator.StructCtx(ctx, query); err != nil {
			validationErr := customErrors.NewValidationErrorWrap(
				err,
				"[searchProductsEndpoint_handler.StructCtx]  query validation failed",
			)
			ep.Logger.Errorf("[searchProductsEndpoint_handler.StructCtx] err: {%v}", validationErr)
			return validationErr
		}

		queryResult, err := mediatr.Send[*queries.SearchProducts, *dtos.SearchProductsResponseDto](
			ctx,
			query,
		)
		if err != nil {
			err = errors.WithMessage(
				err,
				"[searchProductsEndpoint_handler.Send] error in sending SearchProducts",
			)
			ep.Logger.Error(fmt.Sprintf("[searchProductsEndpoint_handler.Send] err: {%v}", err))
			return err
		}

		return c.JSON(http.StatusOK, queryResult)
	}
}
