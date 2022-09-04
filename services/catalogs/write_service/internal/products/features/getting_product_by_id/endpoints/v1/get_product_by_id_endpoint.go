package v1

import (
	"fmt"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_product_by_id/queries/v1"
	"github.com/pkg/errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/delivery"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_product_by_id/dtos"
)

type getProductByIdEndpoint struct {
	*delivery.ProductEndpointBase
}

func NewGetProductByIdEndpoint(productEndpointBase *delivery.ProductEndpointBase) *getProductByIdEndpoint {
	return &getProductByIdEndpoint{productEndpointBase}
}

func (ep *getProductByIdEndpoint) MapRoute() {
	ep.ProductsGroup.GET("/:id", ep.handler())
}

// GetProductByID
// @Tags Products
// @Summary Get product by id
// @Description Get product by id
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} dtos.GetProductByIdResponseDto
// @Router /api/v1/products/{id} [get]
func (ep *getProductByIdEndpoint) handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		ep.Metrics.GetProductByIdHttpRequests.Inc()
		ctx, span := tracing.StartHttpServerTracerSpan(c, "getProductByIdEndpoint.handler")
		defer span.Finish()

		request := &dtos.GetProductByIdRequestDto{}
		if err := c.Bind(request); err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(err, "[getProductByIdEndpoint_handler.Bind] error in the binding request")
			ep.Log.Errorf(fmt.Sprintf("[getProductByIdEndpoint_handler.Bind] err: %v", tracing.TraceWithErr(span, badRequestErr)))
			return badRequestErr
		}

		query := v1.NewGetProductByIdQuery(request.ProductId)
		if err := ep.Validator.StructCtx(ctx, query); err != nil {
			validationErr := customErrors.NewValidationErrorWrap(err, "[getProductByIdEndpoint_handler.StructCtx]  query validation failed")
			ep.Log.Errorf("[getProductByIdEndpoint_handler.StructCtx] err: {%v}", tracing.TraceWithErr(span, validationErr))
			return validationErr
		}

		queryResult, err := mediatr.Send[*v1.GetProductByIdQuery, *dtos.GetProductByIdResponseDto](ctx, query)

		if err != nil {
			err = errors.WithMessage(err, "[getProductByIdEndpoint_handler.Send] error in sending GetProductByIdQuery")
			ep.Log.Errorw(fmt.Sprintf("[getProductByIdEndpoint_handler.Send] id: {%s}, err: {%v}", query.ProductID, tracing.TraceWithErr(span, err)), logger.Fields{"productId": query.ProductID})
			return err
		}

		return c.JSON(http.StatusOK, queryResult)
	}
}
