package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/delivery"
	creatingOrderv1 "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/creating_order/commands/v1"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/creating_order/dtos"
	"net/http"
	"time"
)

type createOrderEndpoint struct {
	*delivery.OrderEndpointBase
}

func NewCreteOrderEndpoint(endpointBase *delivery.OrderEndpointBase) *createOrderEndpoint {
	return &createOrderEndpoint{endpointBase}
}

func (ep *createOrderEndpoint) MapRoute() {
	ep.OrdersGroup.POST("", ep.createOrder())
}

// Create Order
// @Tags Orders
// @Summary Create order
// @Description Create new order
// @Accept json
// @Produce json
// @Param CreateOrderRequestDto body dtos.CreateOrderRequestDto true "Order data"
// @Success 201 {object} dtos.CreateOrderResponseDto
// @Router /api/v1/orders [post]
func (ep *createOrderEndpoint) createOrder() echo.HandlerFunc {
	return func(c echo.Context) error {
		ep.Metrics.CreateOrderHttpRequests.Inc()
		ctx, span := tracing.StartHttpServerTracerSpan(c, "createOrderEndpoint.createOrder")
		defer span.Finish()

		request := &dtos.CreateOrderRequestDto{}
		if err := c.Bind(request); err != nil {
			ep.Log.WarnMsg("Bind", err)
			tracing.TraceErr(span, err)
			return err
		}

		if err := ep.Validator.StructCtx(ctx, request); err != nil {
			ep.Log.Errorf("(validate) err: {%v}", err)
			tracing.TraceErr(span, err)
			return err
		}

		command := creatingOrderv1.NewCreateOrderCommand(request.ShopItems, request.AccountEmail, request.DeliveryAddress, time.Time(request.DeliveryTime))
		result, err := mediatr.Send[*creatingOrderv1.CreateOrderCommand, *dtos.CreateOrderResponseDto](ctx, command)

		if err != nil {
			ep.Log.Errorf("(CreateOrder.Handle) id: {%s}, err: {%v}", command.OrderID, err)
			tracing.TraceErr(span, err)
			return err
		}

		ep.Log.Infof("(order created) id: {%s}", command.OrderID)
		return c.JSON(http.StatusCreated, result)
	}
}
