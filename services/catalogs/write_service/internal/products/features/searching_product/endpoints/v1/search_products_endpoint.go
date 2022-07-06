package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/delivery"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/searching_product"
	"net/http"
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

// SearchProducts
// @Tags Products
// @Summary Search products
// @Description Search products
// @Accept json
// @Produce json
// @Param searchProductsRequestDto query searching_product.SearchProductsRequestDto false "SearchProductsRequestDto"
// @Success 200 {object} searching_product.SearchProductsResponseDto
// @Router /products/search [get]
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

		request := &searching_product.SearchProductsRequestDto{ListQuery: listQuery}

		// https://echo.labstack.com/guide/binding/
		if err := c.Bind(request); err != nil {
			ep.Log.WarnMsg("Bind", err)
			tracing.TraceErr(span, err)
			return err
		}

		query := searching_product.SearchProducts{SearchText: request.SearchText, ListQuery: request.ListQuery}

		if err := ep.Validator.StructCtx(ctx, query); err != nil {
			ep.Log.Errorf("(validate) err: {%v}", err)
			tracing.TraceErr(span, err)
			return err
		}

		queryResult, err := ep.ProductMediator.Send(ctx, query)

		if err != nil {
			ep.Log.WarnMsg("SearchProducts", err)
			tracing.TraceErr(span, err)
			return err
		}

		response, ok := queryResult.(*searching_product.SearchProductsResponseDto)
		err = utils.CheckType(ok)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		return c.JSON(http.StatusOK, response)
	}
}
