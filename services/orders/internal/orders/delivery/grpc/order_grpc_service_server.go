package grpc

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/configurations/infrastructure"
)

type OrderGrpcServiceServer struct {
	*infrastructure.InfrastructureConfiguration
}

func NewOrderGrpcService(infra *infrastructure.InfrastructureConfiguration) *OrderGrpcServiceServer {
	return &OrderGrpcServiceServer{InfrastructureConfiguration: infra}
}

//
//func (s *OrderGrpcServiceServer) CreateOrder(ctx context.Context, req *orders_service.CreateOrderReq) (*orders_service.CreateOrderRes, error) {
//	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "OrderGrpcServiceServer.CreateOrder")
//	defer span.Finish()
//	span.LogFields(log.String("req", req.String()))
//	s.Metrics.CreateOrderGrpcRequests.Inc()
//
//	shopItemsDtos, err := mapper.Map[[]*dtos.ShopItemDto](req.GetShopItems())
//	if err != nil {
//		return nil, err
//	}
//
//	command := v1.NewCreateOrderCommand(shopItemsDtos, req.GetAccountEmail(), req.GetDeliveryAddress())
//	if err := s.Validator.StructCtx(ctx, command); err != nil {
//		s.Log.Errorf("(validate) aggregateID: {%s}, err: {%v}", command.OrderID, err)
//		tracing.TraceErr(span, err)
//		return nil, s.errResponse(err)
//	}
//
//	if _, err := mediatr.Send[*mediatr.Unit](ctx, command); err != nil {
//		s.Log.Errorf("(CreateOrder.Handle) orderID: {%s}, err: {%v}", command.OrderID, err)
//		return nil, s.errResponse(err)
//	}
//
//	s.Log.Infof("(created order): orderID: {%s}", command.OrderID)
//	return &orders_service.CreateOrderRes{AggregateID: command.OrderID.String()}, nil
//}
//
//func (s *OrderGrpcServiceServer) PayOrder(ctx context.Context, req *orders_service.PayOrderReq) (*orders_service.PayOrderRes, error) {
//	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "orderGrpcService.PayOrder")
//	defer span.Finish()
//	span.LogFields(log.String("req", req.String()))
//	s.Metrics.PayOrderGrpcRequests.Inc()
//
//	paymentId, err := uuid.FromString(req.GetPayment().GetID())
//	if err != nil {
//		return nil, err
//	}
//
//	orderId, err := uuid.FromString(req.GetAggregateID())
//
//	payment := models.Payment{PaymentID: paymentId, OrderID: orderId, Timestamp: time.Now()}
//	command := v1.NewPayOrderCommand(payment, req.GetAggregateID())
//	if err := s.v.StructCtx(ctx, command); err != nil {
//		s.Log.Errorf("(validate) err: {%v}", err)
//		tracing.TraceErr(span, err)
//		return nil, s.errResponse(err)
//	}
//
//	if err := s.os.Commands.OrderPaid.Handle(ctx, command); err != nil {
//		s.Log.Errorf("(OrderPaid.Handle) orderID: {%s}, err: {%v}", req.GetAggregateID(), err)
//		return nil, s.errResponse(err)
//	}
//
//	s.Log.Infof("(paid order): orderID: {%s}", req.GetAggregateID())
//	return &orders_service.PayOrderRes{AggregateID: req.GetAggregateID()}, nil
//}
//
//func (s *OrderGrpcServiceServer) SubmitOrder(ctx context.Context, req *orders_service.SubmitOrderReq) (*orders_service.SubmitOrderRes, error) {
//	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "orderGrpcService.SubmitOrder")
//	defer span.Finish()
//	span.LogFields(log.String("req", req.String()))
//	s.Metrics.SubmitOrderGrpcRequests.Inc()
//
//	command := v1.NewSubmitOrderCommand(req.GetAggregateID())
//	if err := s.v.StructCtx(ctx, command); err != nil {
//		s.Log.Errorf("(validate) err: {%v}", err)
//		tracing.TraceErr(span, err)
//		return nil, s.errResponse(err)
//	}
//
//	if err := s.os.Commands.SubmitOrder.Handle(ctx, command); err != nil {
//		s.Log.Errorf("(SubmitOrder.Handle) orderID: {%s}, err: {%v}", req.GetAggregateID(), err)
//		return nil, s.errResponse(err)
//	}
//
//	s.Log.Infof("(submitted order): orderID: {%s}", req.GetAggregateID())
//	return &orders_service.SubmitOrderRes{AggregateID: req.GetAggregateID()}, nil
//}
//
//func (s *OrderGrpcServiceServer) GetOrderByID(ctx context.Context, req *orders_service.GetOrderByIDReq) (*orders_service.GetOrderByIDRes, error) {
//	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "orderGrpcService.GetOrderByID")
//	defer span.Finish()
//	span.LogFields(log.String("req", req.String()))
//	s.Metrics.GetOrderByIdGrpcRequests.Inc()
//
//	query := queries.NewGetOrderByIDQuery(req.GetAggregateID())
//	if err := s.v.StructCtx(ctx, query); err != nil {
//		s.Log.Errorf("(validate) err: {%v}", err)
//		tracing.TraceErr(span, err)
//		return nil, s.errResponse(err)
//	}
//
//	orderProjection, err := s.os.Queries.GetOrderByID.Handle(ctx, query)
//	if err != nil {
//		s.Log.Errorf("(GetOrderByID.Handle) orderID: {%s}, err: {%v}", req.GetAggregateID(), err)
//		return nil, s.errResponse(err)
//	}
//
//	s.Log.Infof("(GetOrderByID) AggregateID: {%s}", req.GetAggregateID())
//	s.Log.Debugf("(GetOrderByID) orderProjection: {%s}", orderProjection.String())
//	return &orders_service.GetOrderByIDRes{Order: models.OrderProjectionToProto(orderProjection)}, nil
//}
//
//func (s *OrderGrpcServiceServer) UpdateShoppingCart(ctx context.Context, req *orders_service.UpdateShoppingCartReq) (*orders_service.UpdateShoppingCartRes, error) {
//	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "orderGrpcService.UpdateShoppingCart")
//	defer span.Finish()
//	span.LogFields(log.String("UpdateShoppingCart req", req.String()))
//	s.Metrics.UpdateOrderGrpcRequests.Inc()
//
//	command := v1.NewUpdateShoppingCartCommand(req.GetAggregateID(), models.ShopItemsFromProto(req.GetShopItems()))
//	if err := s.v.StructCtx(ctx, command); err != nil {
//		s.Log.Errorf("(validate) err: {%v}", err)
//		tracing.TraceErr(span, err)
//		return nil, s.errResponse(err)
//	}
//
//	if err := s.os.Commands.UpdateOrder.Handle(ctx, command); err != nil {
//		s.Log.Errorf("(UpdateShoppingCart.Handle) orderID: {%s}, err: {%v}", req.GetAggregateID(), err)
//		return nil, s.errResponse(err)
//	}
//
//	s.Log.Infof("(UpdateShoppingCart): AggregateID: {%s}", req.GetAggregateID())
//	return &orders_service.UpdateShoppingCartRes{}, nil
//}
//
//func (s *OrderGrpcServiceServer) CancelOrder(ctx context.Context, req *orders_service.CancelOrderReq) (*orders_service.CancelOrderRes, error) {
//	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "orderGrpcService.CancelOrder")
//	defer span.Finish()
//	span.LogFields(log.String("CancelOrder req", req.String()))
//	s.Metrics.CancelOrderGrpcRequests.Inc()
//
//	command := v1.NewCancelOrderCommand(req.GetAggregateID(), req.GetCancelReason())
//	if err := s.v.StructCtx(ctx, command); err != nil {
//		s.Log.Errorf("(validate) err: {%v}", err)
//		tracing.TraceErr(span, err)
//		return nil, s.errResponse(err)
//	}
//
//	if err := s.os.Commands.CancelOrder.Handle(ctx, command); err != nil {
//		s.Log.Errorf("(CancelOrder.Handle) orderID: {%s}, err: {%v}", req.GetAggregateID(), err)
//		return nil, s.errResponse(err)
//	}
//
//	s.Log.Infof("(CancelOrder): AggregateID: {%s}", req.GetAggregateID())
//	return &orders_service.CancelOrderRes{}, nil
//}
//
//func (s *OrderGrpcServiceServer) CompleteOrder(ctx context.Context, req *orders_service.CompleteOrderReq) (*orders_service.CompleteOrderRes, error) {
//	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "orderGrpcService.CompleteOrder")
//	defer span.Finish()
//	span.LogFields(log.String("CompleteOrder req", req.String()))
//	s.Metrics.CompleteOrderGrpcRequests.Inc()
//
//	command := v1.NewCompleteOrderCommand(req.GetAggregateID(), time.Now())
//	if err := s.v.StructCtx(ctx, command); err != nil {
//		s.Log.Errorf("(validate) err: {%v}", err)
//		tracing.TraceErr(span, err)
//		return nil, s.errResponse(err)
//	}
//
//	if err := s.os.Commands.CompleteOrder.Handle(ctx, command); err != nil {
//		s.Log.Errorf("(CompleteOrder.Handle) orderID: {%s}, err: {%v}", req.GetAggregateID(), err)
//		return nil, s.errResponse(err)
//	}
//
//	s.Log.Infof("(CompleteOrder): AggregateID: {%s}", req.GetAggregateID())
//	return &orders_service.CompleteOrderRes{}, nil
//}
//
//func (s *OrderGrpcServiceServer) ChangeDeliveryAddress(ctx context.Context, req *orders_service.ChangeDeliveryAddressReq) (*orders_service.ChangeDeliveryAddressRes, error) {
//	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "orderGrpcService.ChangeDeliveryAddress")
//	defer span.Finish()
//	span.LogFields(log.String("ChangeDeliveryAddress req", req.String()))
//	s.Metrics.ChangeAddressOrderGrpcRequests.Inc()
//
//	command := v1.NewChangeDeliveryAddressCommand(req.GetAggregateID(), req.GetDeliveryAddress())
//	if err := s.v.StructCtx(ctx, command); err != nil {
//		s.Log.Errorf("(validate) err: {%v}", err)
//		tracing.TraceErr(span, err)
//		return nil, s.errResponse(err)
//	}
//
//	if err := s.os.Commands.ChangeOrderDeliveryAddress.Handle(ctx, command); err != nil {
//		s.Log.Errorf("(ChangeOrderDeliveryAddress.Handle) orderID: {%s}, err: {%v}", req.GetAggregateID(), err)
//		return nil, s.errResponse(err)
//	}
//
//	s.Log.Infof("(ChangeDeliveryAddress): AggregateID: {%s}", req.GetAggregateID())
//	return &orders_service.ChangeDeliveryAddressRes{}, nil
//}
//
//func (s *OrderGrpcServiceServer) Search(ctx context.Context, req *orders_service.SearchReq) (*orders_service.SearchRes, error) {
//	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "orderGrpcService.Search")
//	defer span.Finish()
//	span.LogFields(log.String("SearchText", req.GetSearchText()), log.Int64("Page", req.GetPage()), log.Int64("Size", req.GetSize()))
//	s.Metrics.SearchOrderGrpcRequests.Inc()
//
//	query := queries.NewSearchOrdersQuery(req.GetSearchText(), utils.NewPaginationQuery(int(req.GetSize()), int(req.GetPage())))
//	if err := s.v.StructCtx(ctx, query); err != nil {
//		s.Log.Errorf("(validate) err: {%v}", err)
//		tracing.TraceErr(span, err)
//		return nil, s.errResponse(err)
//	}
//
//	searchResult, err := s.os.Queries.SearchOrders.Handle(ctx, query)
//	if err != nil {
//		s.Log.Errorf("(SearchOrders.Handle) text: {%s}, err: {%v}", req.GetSearchText(), err)
//		return nil, s.errResponse(err)
//	}
//
//	s.Log.Infof("(Search result): searchText: {%s}, pagination: {%+v}", req.GetSearchText(), searchResult.Pagination)
//	return mappers.SearchResponseToProto(searchResult), nil
//}
//
//func (s *OrderGrpcServiceServer) errResponse(err error) error {
//	return grpcErrors.ErrResponse(err)
//}
