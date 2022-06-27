package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/constants"
	httpErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts/repositories"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_product_by_id"
	shared_configurations "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/configurations"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

type getProductByIdEndpoint struct {
	mediator          *mediatr.Mediator
	productRepository repositories.ProductRepository
	productsGroup     *echo.Group
	infrastructure    *shared_configurations.Infrastructure
}

func NewGetProductByIdEndpoint(infra *shared_configurations.Infrastructure, mediator *mediatr.Mediator, productsGroup *echo.Group, productRepository repositories.ProductRepository) *getProductByIdEndpoint {
	return &getProductByIdEndpoint{mediator: mediator, productRepository: productRepository, productsGroup: productsGroup, infrastructure: infra}
}

func (ep *getProductByIdEndpoint) MapRoute() {
	ep.productsGroup.GET("/:id", ep.getProductByID())
}

// GetProductByID
// @Tags Products
// @Summary Get product
// @Description Get product by id
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} getting_product_by_id.GetProductByIdResponseDto
// @Router /products/{id} [get]
func (ep *getProductByIdEndpoint) getProductByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		ep.infrastructure.Metrics.GetProductByIdHttpRequests.Inc()

		ctx, span := tracing.StartHttpServerTracerSpan(c, "productsHandlers.getProductByID")
		defer span.Finish()

		productUUID, err := uuid.FromString(c.Param(constants.ID))
		if err != nil {
			ep.infrastructure.Log.WarnMsg("uuid.FromString", err)
			ep.infrastructure.TraceErr(span, err)
			return httpErrors.ErrorResponse(err, ep.infrastructure.Cfg.Http.DebugErrorsResponse)
		}

		query := getting_product_by_id.NewGetProductById(productUUID)
		queryResult, err := ep.mediator.Send(ctx, query)

		if err != nil {
			ep.infrastructure.Log.WarnMsg("GetProductById", err)
			ep.infrastructure.Metrics.ErrorHttpRequests.Inc()
			return httpErrors.ErrorResponse(err, ep.infrastructure.Cfg.Http.DebugErrorsResponse)
		}

		response, ok := queryResult.(*getting_product_by_id.GetProductByIdResponseDto)
		err = utils.CheckType(ok)
		if err != nil {
			return httpErrors.ErrorResponse(err, ep.infrastructure.Cfg.Http.DebugErrorsResponse)
		}

		ep.infrastructure.Metrics.SuccessHttpRequests.Inc()
		return c.JSON(http.StatusOK, response)
	}
}
