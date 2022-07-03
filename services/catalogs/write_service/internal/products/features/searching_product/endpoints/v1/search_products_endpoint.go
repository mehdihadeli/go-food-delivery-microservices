package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts/repositories"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/searching_product"
	shared_configurations "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/configurations"
	"net/http"
)

type searchProductsEndpoint struct {
	mediator          *mediatr.Mediator
	productRepository repositories.ProductRepository
	productsGroup     *echo.Group
	infrastructure    *shared_configurations.Infrastructure
}

func NewSearchProductsEndpoint(infra *shared_configurations.Infrastructure, mediator *mediatr.Mediator, productsGroup *echo.Group, productRepository repositories.ProductRepository) *searchProductsEndpoint {
	return &searchProductsEndpoint{mediator: mediator, productRepository: productRepository, productsGroup: productsGroup, infrastructure: infra}
}

func (ep *searchProductsEndpoint) MapRoute() {
	ep.productsGroup.GET("/search", ep.searchProducts())
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
		ep.infrastructure.Metrics.GetProductByIdHttpRequests.Inc()

		ctx, span := tracing.StartHttpServerTracerSpan(c, "searchProductsEndpoint.searchProducts")
		defer span.Finish()

		listQuery, err := utils.GetListQueryFromCtx(c)

		if err != nil {
			utils.LogResponseError(c, ep.infrastructure.Log, err)
			return err
		}

		request := &searching_product.SearchProductsRequestDto{ListQuery: listQuery}

		// https://echo.labstack.com/guide/binding/
		if err := c.Bind(request); err != nil {
			ep.infrastructure.Log.WarnMsg("Bind", err)
			ep.infrastructure.TraceErr(span, err)
			return err
		}

		if err := ep.infrastructure.Validator.StructCtx(ctx, request); err != nil {
			ep.infrastructure.Log.Errorf("(validate) err: {%v}", err)
			tracing.TraceErr(span, err)
			return err
		}

		query := searching_product.SearchProducts{SearchText: request.SearchText, ListQuery: request.ListQuery}

		queryResult, err := ep.mediator.Send(ctx, query)

		if err != nil {
			ep.infrastructure.Log.WarnMsg("SearchProducts", err)
			ep.infrastructure.Metrics.ErrorHttpRequests.Inc()
			return err
		}

		response, ok := queryResult.(*searching_product.SearchProductsResponseDto)
		err = utils.CheckType(ok)
		if err != nil {
			return err
		}

		ep.infrastructure.Metrics.SuccessHttpRequests.Inc()
		return c.JSON(http.StatusOK, response)
	}
}
