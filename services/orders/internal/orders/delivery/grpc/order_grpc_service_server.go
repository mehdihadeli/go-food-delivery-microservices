package grpc

import (
	"context"
	"emperror.dev/errors"
	"fmt"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc/grpcErrors"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mapper"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	grpcOrderService "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/contracts/proto/service_clients"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/dtos"
	creatingOrderCommandV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/creating_order/commands/v1"
	orderDtos "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/creating_order/dtos"
	gettingOrderByIdDtos "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/getting_order_by_id/dtos"
	gettingOrderByIdQueryV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/getting_order_by_id/queries/v1"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/models/orders/aggregate"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/configurations/infrastructure"
	"github.com/opentracing/opentracing-go/log"
	uuid "github.com/satori/go.uuid"
)

type OrderGrpcServiceServer struct {
	*infrastructure.InfrastructureConfiguration
}

func NewOrderGrpcService(infra *infrastructure.InfrastructureConfiguration) *OrderGrpcServiceServer {
	return &OrderGrpcServiceServer{InfrastructureConfiguration: infra}
}

func (o OrderGrpcServiceServer) CreateOrder(ctx context.Context, req *grpcOrderService.CreateOrderReq) (*grpcOrderService.CreateOrderRes, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrderGrpcServiceServer.CreateOrder")
	span.LogFields(log.Object("Request", req))
	o.Metrics.CreateOrderGrpcRequests.Inc()
	defer span.Finish()

	shopItemsDtos, err := mapper.Map[[]*dtos.ShopItemDto](req.GetShopItems())
	if err != nil {
		return nil, err
	}

	command := creatingOrderCommandV1.NewCreateOrderCommand(shopItemsDtos, req.AccountEmail, req.DeliveryAddress, req.DeliveryTime.AsTime())
	if err := o.Validator.StructCtx(ctx, command); err != nil {
		validationErr := customErrors.NewValidationErrorWrap(err, "[OrderGrpcServiceServer_CreateOrder.StructCtx] command validation failed")
		o.Log.Errorf(fmt.Sprintf("[OrderGrpcServiceServer_CreateOrder.StructCtx] err: %v", tracing.TraceWithErr(span, validationErr)))
		return nil, grpcErrors.ErrGrpcResponse(validationErr)
	}

	result, err := mediatr.Send[*creatingOrderCommandV1.CreateOrderCommand, *orderDtos.CreateOrderResponseDto](ctx, command)

	if err != nil {
		err = errors.WithMessage(err, "[ProductGrpcServiceServer_CreateOrder.Send] error in sending CreateOrderCommand")
		o.Log.Errorw(fmt.Sprintf("[ProductGrpcServiceServer_CreateOrder.Send] id: {%s}, err: %v", command.OrderID, tracing.TraceWithErr(span, err)), logger.Fields{"OrderId": command.OrderID})
		return nil, grpcErrors.ErrGrpcResponse(err)
	}

	o.Metrics.SuccessGrpcRequests.Inc()
	return &grpcOrderService.CreateOrderRes{OrderId: result.OrderID.String()}, nil
}

func (o OrderGrpcServiceServer) GetOrderByID(ctx context.Context, req *grpcOrderService.GetOrderByIDReq) (*grpcOrderService.GetOrderByIDRes, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrderGrpcServiceServer.GetOrderByID")
	span.LogFields(log.Object("Request", req))
	o.Metrics.GetOrderByIdGrpcRequests.Inc()
	defer span.Finish()

	orderIdUUID, err := uuid.FromString(req.OrderId)
	if err != nil {
		badRequestErr := customErrors.NewBadRequestErrorWrap(err, "[OrderGrpcServiceServer_GetOrderByID.uuid.FromString] error in converting uuid")
		o.Log.Errorf(fmt.Sprintf("[OrderGrpcServiceServer_GetOrderByID.uuid.FromString] err: %v", tracing.TraceWithErr(span, badRequestErr)))
		return nil, grpcErrors.ErrGrpcResponse(badRequestErr)
	}

	query := gettingOrderByIdQueryV1.NewGetOrderByIdQuery(orderIdUUID)
	if err := o.Validator.StructCtx(ctx, query); err != nil {
		validationErr := customErrors.NewValidationErrorWrap(err, "[OrderGrpcServiceServer_GetOrderByID.StructCtx] query validation failed")
		o.Log.Errorf(fmt.Sprintf("[OrderGrpcServiceServer_GetOrderByID.StructCtx] err: %v", tracing.TraceWithErr(span, validationErr)))
		return nil, grpcErrors.ErrGrpcResponse(validationErr)
	}

	queryResult, err := mediatr.Send[*gettingOrderByIdQueryV1.GetOrderByIdQuery, *gettingOrderByIdDtos.GetOrderByIdResponseDto](ctx, query)
	if err != nil {
		err = errors.WithMessage(err, "[OrderGrpcServiceServer_GetOrderByID.Send] error in sending GetOrderByIdQuery")
		o.Log.Errorw(fmt.Sprintf("[OrderGrpcServiceServer_GetOrderByID.Send] id: {%s}, err: %v", query.OrderId, tracing.TraceWithErr(span, err)), logger.Fields{"OrderId": query.OrderId})
		return nil, grpcErrors.ErrGrpcResponse(err)
	}

	order, err := mapper.Map[*aggregate.Order](queryResult.Order)
	if err != nil {
		err = errors.WithMessage(err, "[OrderGrpcServiceServer_GetOrderByID.Map] error in mapping order")
		return nil, grpcErrors.ErrGrpcResponse(tracing.TraceWithErr(span, err))
	}

	ord, err := mapper.Map[*grpcOrderService.Order](order)
	if err != nil {
		err = errors.WithMessage(err, "[ProductGrpcServiceServer_GetOrderByID.Map] error in mapping order")
		return nil, grpcErrors.ErrGrpcResponse(tracing.TraceWithErr(span, err))
	}

	o.Metrics.SuccessGrpcRequests.Inc()

	return &grpcOrderService.GetOrderByIDRes{Order: ord}, nil
}

func (o OrderGrpcServiceServer) SubmitOrder(ctx context.Context, req *grpcOrderService.SubmitOrderReq) (*grpcOrderService.SubmitOrderRes, error) {
	//TODO implement me
	panic("implement me")
}

func (o OrderGrpcServiceServer) UpdateShoppingCart(ctx context.Context, req *grpcOrderService.UpdateShoppingCartReq) (*grpcOrderService.UpdateShoppingCartRes, error) {
	//TODO implement me
	panic("implement me")
}

func (o OrderGrpcServiceServer) GetOrders(ctx context.Context, req *grpcOrderService.GetOrdersReq) (*grpcOrderService.GetOrdersRes, error) {
	//TODO implement me
	panic("implement me")
}
