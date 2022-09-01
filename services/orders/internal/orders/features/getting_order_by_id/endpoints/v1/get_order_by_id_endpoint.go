package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/delivery"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/getting_order_by_id/dtos"
	v1 "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/getting_order_by_id/queries/v1"
	"net/http"
)

type getOrderByIdEndpoint struct {
	*delivery.OrderEndpointBase
}

func NewGetOrderByIdEndpoint(productEndpointBase *delivery.OrderEndpointBase) *getOrderByIdEndpoint {
	return &getOrderByIdEndpoint{productEndpointBase}
}

func (ep *getOrderByIdEndpoint) MapRoute() {
	ep.OrdersGroup.GET("/:id", ep.getOrderByID())
}

// Get Order By ID
// @Tags Orders
// @Summary Get order by id
// @Description Get order by id
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} dtos.GetOrderByIdResponseDto
// @Router /api/v1/orders/{id} [get]
func (ep *getOrderByIdEndpoint) getOrderByID() echo.HandlerFunc {
	return func(c echo.Context) error {

		ep.Metrics.GetOrderByIdHttpRequests.Inc()
		ctx, span := tracing.StartHttpServerTracerSpan(c, "getOrderByIdEndpoint.getOrderByID")
		defer span.Finish()

		request := &dtos.GetOrderByIdRequestDto{}
		if err := c.Bind(request); err != nil {
			ep.Log.Errorf("(Bind) err: %v", tracing.TraceWithErr(span, err))
			return err
		}

		query := &v1.GetOrderByIdQuery{OrderId: request.OrderId}

		if err := ep.Validator.StructCtx(ctx, query); err != nil {
			ep.Log.Errorf("(validate) err: %v", tracing.TraceWithErr(span, err))
			return err
		}

		queryResult, err := mediatr.Send[*v1.GetOrderByIdQuery, *dtos.GetOrderByIdResponseDto](ctx, query)

		if err != nil {
			ep.Log.Errorf("(GetOrderById.Handle) id: %s, err: %v", query.OrderId, tracing.TraceWithErr(span, err))
			return err
		}

		return c.JSON(http.StatusOK, queryResult)
	}
}
