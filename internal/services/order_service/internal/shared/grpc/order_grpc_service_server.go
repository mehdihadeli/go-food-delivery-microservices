package grpc

import (
	"context"
	"fmt"

	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/mapper"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing"
	attribute2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing/attribute"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils"
	dtosV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/dtos/v1"
	createOrderCommandV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/features/creating_order/v1/commands"
	createOrderDtosV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/features/creating_order/v1/dtos"
	getOrderByIdDtosV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/features/getting_order_by_id/v1/dtos"
	getOrderByIdQueryV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/features/getting_order_by_id/v1/queries"
	getOrdersDtosV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/features/getting_orders/v1/dtos"
	getOrdersQueryV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/features/getting_orders/v1/queries"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/shared/contracts"
	grpcOrderService "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/shared/grpc/genproto"

	"emperror.dev/errors"
	"github.com/go-playground/validator"
	"github.com/mehdihadeli/go-mediatr"
	uuid "github.com/satori/go.uuid"
	"go.opentelemetry.io/otel/attribute"
	api "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type OrderGrpcServiceServer struct {
	ordersMetrics *contracts.OrdersMetrics
	logger        logger.Logger
	validator     *validator.Validate
}

var grpcMetricsAttr = api.WithAttributes(
	attribute.Key("MetricsType").String("Grpc"),
)

func NewOrderGrpcService(
	logger logger.Logger,
	validator *validator.Validate,
	ordersMetrics *contracts.OrdersMetrics,
) *OrderGrpcServiceServer {
	return &OrderGrpcServiceServer{
		ordersMetrics: ordersMetrics,
		logger:        logger,
		validator:     validator,
	}
}

func (o OrderGrpcServiceServer) CreateOrder(
	ctx context.Context,
	req *grpcOrderService.CreateOrderReq,
) (*grpcOrderService.CreateOrderRes, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute2.Object("Request", req))
	o.ordersMetrics.CreateOrderGrpcRequests.Add(ctx, 1, grpcMetricsAttr)

	shopItemsDtos, err := mapper.Map[[]*dtosV1.ShopItemDto](req.GetShopItems())
	if err != nil {
		return nil, err
	}

	command, err := createOrderCommandV1.NewCreateOrder(
		shopItemsDtos,
		req.AccountEmail,
		req.DeliveryAddress,
		req.DeliveryTime.AsTime(),
	)
	if err != nil {
		validationErr := customErrors.NewValidationErrorWrap(
			err,
			"[OrderGrpcServiceServer_CreateOrder.StructCtx] command validation failed",
		)
		o.logger.Errorf(
			fmt.Sprintf("[OrderGrpcServiceServer_CreateOrder.StructCtx] err: %v", validationErr),
		)
		return nil, validationErr
	}

	result, err := mediatr.Send[*createOrderCommandV1.CreateOrder, *createOrderDtosV1.CreateOrderResponseDto](
		ctx,
		command,
	)
	if err != nil {
		err = errors.WithMessage(
			err,
			"[ProductGrpcServiceServer_CreateOrder.Send] error in sending CreateOrder",
		)
		o.logger.Errorw(
			fmt.Sprintf(
				"[ProductGrpcServiceServer_CreateOrder.Send] id: {%s}, err: %v",
				command.OrderId,
				err,
			),
			logger.Fields{"Id": command.OrderId},
		)
		return nil, err
	}

	return &grpcOrderService.CreateOrderRes{OrderId: result.OrderId.String()}, nil
}

func (o OrderGrpcServiceServer) GetOrderByID(
	ctx context.Context,
	req *grpcOrderService.GetOrderByIDReq,
) (*grpcOrderService.GetOrderByIDRes, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute2.Object("Request", req))
	o.ordersMetrics.GetOrderByIdGrpcRequests.Add(ctx, 1, grpcMetricsAttr)

	orderIdUUID, err := uuid.FromString(req.Id)
	if err != nil {
		badRequestErr := customErrors.NewBadRequestErrorWrap(
			err,
			"[OrderGrpcServiceServer_GetOrderByID.uuid.FromString] error in converting uuid",
		)
		o.logger.Errorf(
			fmt.Sprintf(
				"[OrderGrpcServiceServer_GetOrderByID.uuid.FromString] err: %v",
				badRequestErr,
			),
		)
		return nil, badRequestErr
	}

	query, err := getOrderByIdQueryV1.NewGetOrderById(orderIdUUID)
	if err != nil {
		validationErr := customErrors.NewValidationErrorWrap(
			err,
			"[OrderGrpcServiceServer_GetOrderByID.StructCtx] query validation failed",
		)
		o.logger.Errorf(
			fmt.Sprintf("[OrderGrpcServiceServer_GetOrderByID.StructCtx] err: %v", validationErr),
		)
		return nil, validationErr
	}

	queryResult, err := mediatr.Send[*getOrderByIdQueryV1.GetOrderById, *getOrderByIdDtosV1.GetOrderByIdResponseDto](
		ctx,
		query,
	)
	if err != nil {
		err = errors.WithMessage(
			err,
			"[OrderGrpcServiceServer_GetOrderByID.Send] error in sending GetOrderById",
		)
		o.logger.Errorw(
			fmt.Sprintf(
				"[OrderGrpcServiceServer_GetOrderByID.Send] id: {%s}, err: %v",
				query.Id,
				err,
			),
			logger.Fields{"Id": query.Id},
		)
		return nil, err
	}

	q := queryResult.Order
	order, err := mapper.Map[*grpcOrderService.OrderReadModel](q)
	if err != nil {
		err = errors.WithMessage(
			err,
			"[OrderGrpcServiceServer_GetOrderByID.Map] error in mapping order",
		)
		return nil, tracing.TraceErrFromContext(ctx, err)
	}

	return &grpcOrderService.GetOrderByIDRes{Order: order}, nil
}

func (o OrderGrpcServiceServer) SubmitOrder(
	ctx context.Context,
	req *grpcOrderService.SubmitOrderReq,
) (*grpcOrderService.SubmitOrderRes, error) {
	return nil, nil
}

func (o OrderGrpcServiceServer) UpdateShoppingCart(
	ctx context.Context,
	req *grpcOrderService.UpdateShoppingCartReq,
) (*grpcOrderService.UpdateShoppingCartRes, error) {
	return nil, nil
}

func (o OrderGrpcServiceServer) GetOrders(
	ctx context.Context,
	req *grpcOrderService.GetOrdersReq,
) (*grpcOrderService.GetOrdersRes, error) {
	o.ordersMetrics.GetOrdersGrpcRequests.Add(ctx, 1, grpcMetricsAttr)
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute2.Object("Request", req))

	query := getOrdersQueryV1.NewGetOrders(
		&utils.ListQuery{Page: int(req.Page), Size: int(req.Size)},
	)

	queryResult, err := mediatr.Send[*getOrdersQueryV1.GetOrders, *getOrdersDtosV1.GetOrdersResponseDto](
		ctx,
		query,
	)
	if err != nil {
		err = errors.WithMessage(
			err,
			"[OrderGrpcServiceServer_GetOrders.Send] error in sending GetOrders",
		)
		o.logger.Error(fmt.Sprintf("[OrderGrpcServiceServer_GetOrders.Send] err: {%v}", err))
		return nil, err
	}

	ordersResponse, err := mapper.Map[*grpcOrderService.GetOrdersRes](queryResult.Orders)
	if err != nil {
		err = errors.WithMessage(
			err,
			"[OrderGrpcServiceServer_GetOrders.Map] error in mapping orders",
		)
		return nil, err
	}

	return ordersResponse, nil
}
