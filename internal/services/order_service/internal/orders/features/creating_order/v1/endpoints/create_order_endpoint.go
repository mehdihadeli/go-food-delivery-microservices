package createOrderV1

import (
	"fmt"
	"net/http"
	"time"

	"emperror.dev/errors"
	"github.com/labstack/echo/v4"
	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/web/route"
	"github.com/mehdihadeli/go-mediatr"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/contracts/params"
	createOrderCommandV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/features/creating_order/v1/commands"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/features/creating_order/v1/dtos"
)

type createOrderEndpoint struct {
	params.OrderRouteParams
}

func NewCreteOrderEndpoint(params params.OrderRouteParams) route.Endpoint {
	return &createOrderEndpoint{OrderRouteParams: params}
}

func (ep *createOrderEndpoint) MapEndpoint() {
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
		ctx := c.Request().Context()
		ep.OrdersMetrics.CreateOrderHttpRequests.Add(ctx, 1)

		request := &dtos.CreateOrderRequestDto{}
		if err := c.Bind(request); err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(
				err,
				"[createOrderEndpoint_handler.Bind] error in the binding request",
			)
			ep.Logger.Errorf(
				fmt.Sprintf("[createOrderEndpoint_handler.Bind] err: %v", badRequestErr),
			)
			return badRequestErr
		}

		command := createOrderCommandV1.NewCreateOrder(
			request.ShopItems,
			request.AccountEmail,
			request.DeliveryAddress,
			time.Time(request.DeliveryTime),
		)
		if err := ep.Validator.StructCtx(ctx, command); err != nil {
			validationErr := customErrors.NewValidationErrorWrap(
				err,
				"[createOrderEndpoint_handler.StructCtx] command validation failed",
			)
			ep.Logger.Errorf(
				fmt.Sprintf("[createOrderEndpoint_handler.StructCtx] err: %v", validationErr),
			)
			return validationErr
		}

		result, err := mediatr.Send[*createOrderCommandV1.CreateOrder, *dtos.CreateOrderResponseDto](
			ctx,
			command,
		)
		if err != nil {
			err = errors.WithMessage(
				err,
				"[createOrderEndpoint_handler.Send] error in sending CreateOrder",
			)
			ep.Logger.Errorw(
				fmt.Sprintf(
					"[createOrderEndpoint_handler.Send] id: {%s}, err: %v",
					command.OrderId,
					err,
				),
				logger.Fields{"Id": command.OrderId},
			)
			return err
		}

		return c.JSON(http.StatusCreated, result)
	}
}
