package endpoints

import (
	"fmt"
	"net/http"

	"emperror.dev/errors"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/contracts/params"

	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/web/route"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/features/getting_orders/v1/dtos"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/features/getting_orders/v1/queries"
)

type getOrdersEndpoint struct {
	params.OrderRouteParams
}

func NewGetOrdersEndpoint(params params.OrderRouteParams) route.Endpoint {
	return &getOrdersEndpoint{OrderRouteParams: params}
}

func (ep *getOrdersEndpoint) MapEndpoint() {
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
		ctx := c.Request().Context()
		ep.OrdersMetrics.GetOrdersHttpRequests.Add(ctx, 1)

		listQuery, err := utils.GetListQueryFromCtx(c)
		if err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(
				err,
				"[getOrdersEndpoint_handler.GetListQueryFromCtx] error in getting data from query string",
			)
			ep.Logger.Errorf(
				fmt.Sprintf(
					"[getOrdersEndpoint_handler.GetListQueryFromCtx] err: %v",
					badRequestErr,
				),
			)
			return err
		}

		request := &dtos.GetOrdersRequestDto{ListQuery: listQuery}
		if err := c.Bind(request); err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(
				err,
				"[getOrdersEndpoint_handler.Bind] error in the binding request",
			)
			ep.Logger.Errorf(fmt.Sprintf("[getOrdersEndpoint_handler.Bind] err: %v", badRequestErr))
			return badRequestErr
		}

		query := queries.NewGetOrders(request.ListQuery)

		queryResult, err := mediatr.Send[*queries.GetOrders, *dtos.GetOrdersResponseDto](ctx, query)
		if err != nil {
			err = errors.WithMessage(
				err,
				"[getOrdersEndpoint_handler.Send] error in sending GetOrders",
			)
			ep.Logger.Error(fmt.Sprintf("[getOrdersEndpoint_handler.Send] err: {%v}", err))
			return err
		}

		return c.JSON(http.StatusOK, queryResult)
	}
}
