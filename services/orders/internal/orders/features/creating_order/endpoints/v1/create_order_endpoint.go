package v1

import (
	"emperror.dev/errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
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
	ep.OrdersGroup.POST("", ep.handler())
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
func (ep *createOrderEndpoint) handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		ep.Metrics.CreateOrderHttpRequests.Inc()
		ctx, span := tracing.StartHttpServerTracerSpan(c, "createOrderEndpoint.createOrder")
		defer span.Finish()

		request := &dtos.CreateOrderRequestDto{}
		if err := c.Bind(request); err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(err, "[createOrderEndpoint_handler.Bind] error in the binding request")
			ep.Log.Errorf(fmt.Sprintf("[createOrderEndpoint_handler.Bind] err: %v", tracing.TraceWithErr(span, badRequestErr)))
			return badRequestErr
		}

		command := creatingOrderv1.NewCreateOrder(request.ShopItems, request.AccountEmail, request.DeliveryAddress, time.Time(request.DeliveryTime))
		if err := ep.Validator.StructCtx(ctx, command); err != nil {
			validationErr := customErrors.NewValidationErrorWrap(err, "[createOrderEndpoint_handler.StructCtx] command validation failed")
			ep.Log.Errorf(fmt.Sprintf("[createOrderEndpoint_handler.StructCtx] err: %v", tracing.TraceWithErr(span, validationErr)))
			return validationErr
		}

		result, err := mediatr.Send[*creatingOrderv1.CreateOrder, *dtos.CreateOrderResponseDto](ctx, command)

		if err != nil {
			err = errors.WithMessage(err, "[createOrderEndpoint_handler.Send] error in sending CreateOrder")
			ep.Log.Errorw(fmt.Sprintf("[createOrderEndpoint_handler.Send] id: {%s}, err: %v", command.OrderId, tracing.TraceWithErr(span, err)), logger.Fields{"Id": command.OrderId})
			return err
		}

		return c.JSON(http.StatusCreated, result)
	}
}
