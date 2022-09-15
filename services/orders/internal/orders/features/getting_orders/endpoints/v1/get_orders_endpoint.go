package v1

import (
	"emperror.dev/errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/delivery"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/getting_orders/dtos"
	v1 "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/getting_orders/queryies/v1"
	"net/http"
)

type getOrdersEndpoint struct {
	*delivery.OrderEndpointBase
}

func NewGetOrdersEndpoint(orderEndpointBase *delivery.OrderEndpointBase) *getOrdersEndpoint {
	return &getOrdersEndpoint{orderEndpointBase}
}

func (ep *getOrdersEndpoint) MapRoute() {
	ep.OrdersGroup.GET("", ep.handler())
}

// GetAllOrders
// @Tags Orders
// @Summary Get all orders
// @Description Get all orders
// @Accept json
// @Produce json
// @Param getOrdersRequestDto query dtos.GetOrdersRequestDto false "GetOrdersRequestDto"
// @Success 200 {object} dtos.GetOrdersResponseDto
// @Router /api/v1/orders [get]
func (ep *getOrdersEndpoint) handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		ep.Metrics.GetOrdersHttpRequests.Inc()
		ctx, span := tracing.StartHttpServerTracerSpan(c, "getOrdersEndpoint.handler")
		defer span.Finish()

		listQuery, err := utils.GetListQueryFromCtx(c)
		if err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(err, "[getOrdersEndpoint_handler.GetListQueryFromCtx] error in getting data from query string")
			ep.Log.Errorf(fmt.Sprintf("[getOrdersEndpoint_handler.GetListQueryFromCtx] err: %v", tracing.TraceWithErr(span, badRequestErr)))
			return err
		}

		request := &dtos.GetOrdersRequestDto{ListQuery: listQuery}
		if err := c.Bind(request); err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(err, "[getOrdersEndpoint_handler.Bind] error in the binding request")
			ep.Log.Errorf(fmt.Sprintf("[getOrdersEndpoint_handler.Bind] err: %v", tracing.TraceWithErr(span, badRequestErr)))
			return badRequestErr
		}

		query := v1.NewGetOrders(request.ListQuery)

		queryResult, err := mediatr.Send[*v1.GetOrders, *dtos.GetOrdersResponseDto](ctx, query)

		if err != nil {
			err = errors.WithMessage(err, "[getOrdersEndpoint_handler.Send] error in sending GetOrders")
			ep.Log.Error(fmt.Sprintf("[getOrdersEndpoint_handler.Send] err: {%v}", tracing.TraceWithErr(span, err)))
			return err
		}

		return c.JSON(http.StatusOK, queryResult)
	}
}
