package v1

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/searching_product/queries/v1"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/delivery"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/searching_product/dtos"
)

type searchProductsEndpoint struct {
	*delivery.ProductEndpointBase
}

func NewSearchProductsEndpoint(productEndpointBase *delivery.ProductEndpointBase) *searchProductsEndpoint {
	return &searchProductsEndpoint{productEndpointBase}
}

func (ep *searchProductsEndpoint) MapRoute() {
	ep.ProductsGroup.GET("/search", ep.searchProducts())
}

// SearchProductsQuery
// @Tags Products
// @Summary Search products
// @Description Search products
// @Accept json
// @Produce json
// @Param searchProductsRequestDto query dtos.SearchProductsRequestDto false "SearchProductsRequestDto"
// @Success 200 {object} dtos.SearchProductsResponseDto
// @Router /api/v1/products/search [get]
func (ep *searchProductsEndpoint) searchProducts() echo.HandlerFunc {
	return func(c echo.Context) error {

		ep.Metrics.SearchProductHttpRequests.Inc()
		ctx, span := tracing.StartHttpServerTracerSpan(c, "searchProductsEndpoint.searchProducts")
		defer span.Finish()

		listQuery, err := utils.GetListQueryFromCtx(c)

		if err != nil {
			tracing.TraceErr(span, err)
			utils.LogResponseError(c, ep.Log, err)
			return err
		}

		request := &dtos.SearchProductsRequestDto{ListQuery: listQuery}

		// https://echo.labstack.com/guide/binding/
		if err := c.Bind(request); err != nil {
			ep.Log.WarnMsg("Bind", err)
			tracing.TraceErr(span, err)
			return err
		}

		query := &v1.SearchProductsQuery{SearchText: request.SearchText, ListQuery: request.ListQuery}

		if err := ep.Validator.StructCtx(ctx, query); err != nil {
			ep.Log.Errorf("(validate) err: {%v}", err)
			tracing.TraceErr(span, err)
			return err
		}

		queryResult, err := mediatr.Send[*v1.SearchProductsQuery, *dtos.SearchProductsResponseDto](ctx, query)

		if err != nil {
			ep.Log.WarnMsg("SearchProductsQuery", err)
			tracing.TraceErr(span, err)
			return err
		}

		return c.JSON(http.StatusOK, queryResult)
	}
}
