package v1

import (
	"context"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/contracts/repositories"
	ordersDto "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/dtos"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/getting_orders/dtos"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type GetOrdersQueryHandler struct {
	log                      logger.Logger
	cfg                      *config.Config
	mongoOrderReadRepository repositories.OrderReadRepository
}

func NewGetOrdersQueryHandler(log logger.Logger, cfg *config.Config, mongoOrderReadRepository repositories.OrderReadRepository) *GetOrdersQueryHandler {
	return &GetOrdersQueryHandler{log: log, cfg: cfg, mongoOrderReadRepository: mongoOrderReadRepository}
}

func (c *GetOrdersQueryHandler) Handle(ctx context.Context, query *GetOrdersQuery) (*dtos.GetOrdersResponseDto, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GetOrdersQueryHandler.Handle")
	span.LogFields(log.Object("Query", query))
	defer span.Finish()

	products, err := c.mongoOrderReadRepository.GetAllOrders(ctx, query.ListQuery)
	if err != nil {
		return nil, tracing.TraceWithErr(span, customErrors.NewApplicationErrorWrap(err, "[GetOrdersQueryHandler_Handle.GetAllOrders] error in getting orders in the repository"))
	}

	listResultDto, err := utils.ListResultToListResultDto[*ordersDto.OrderReadDto](products)
	if err != nil {
		return nil, tracing.TraceWithErr(span, customErrors.NewApplicationErrorWrap(err, "[GetOrdersQueryHandler_Handle.ListResultToListResultDto] error in the mapping ListResultToListResultDto"))
	}

	c.log.Info("[GetOrdersQueryHandler.Handle] orders fetched")

	return &dtos.GetOrdersResponseDto{Orders: listResultDto}, nil
}
