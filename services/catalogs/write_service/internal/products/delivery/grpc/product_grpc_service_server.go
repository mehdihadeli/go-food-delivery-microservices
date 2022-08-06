package grpc

import (
	"context"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mapper"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	productsService "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts/proto/service_clients"
	creatingProductV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/creating_product/commands/v1"
	gettingProductByIdV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_product_by_id/queries/v1"
	updatingProductV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/updating_product/commands/v1"

	gettingProductByIdDtos "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_product_by_id/dtos"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/models"

	creatingProductDtos "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/creating_product/dtos"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/configurations/infrastructure"
	"github.com/opentracing/opentracing-go/log"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductGrpcServiceServer struct {
	*infrastructure.InfrastructureConfiguration
	// Ref:https://github.com/grpc/grpc-go/issues/3794#issuecomment-720599532
	// product_service_client.UnimplementedProductsServiceServer
}

func NewProductGrpcService(infra *infrastructure.InfrastructureConfiguration) *ProductGrpcServiceServer {
	return &ProductGrpcServiceServer{InfrastructureConfiguration: infra}
}

func (s *ProductGrpcServiceServer) CreateProduct(ctx context.Context, req *productsService.CreateProductReq) (*productsService.CreateProductRes, error) {

	s.Metrics.CreateProductGrpcRequests.Inc()
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "ProductGrpcServiceServer.CreateProductCommand")
	span.LogFields(log.String("req", req.String()))
	defer span.Finish()

	command := creatingProductV1.NewCreateProductCommand(req.GetName(), req.GetDescription(), req.GetPrice())

	if err := s.Validator.StructCtx(ctx, command); err != nil {
		s.Log.Errorf("(validate) err: {%v}", err)
		tracing.TraceErr(span, err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	result, err := mediatr.Send[*creatingProductV1.CreateProductCommand, *creatingProductDtos.CreateProductResponseDto](ctx, command)
	if err != nil {
		s.Log.Errorf("(CreateProductCommand.Handle) productId: {%s}, err: {%v}", command.ProductID, err)
		tracing.TraceErr(span, err)
		return nil, s.errResponse(codes.Internal, err)
	}

	s.Log.Infof("(product created) productId: {%s}", command.ProductID)
	s.Metrics.SuccessGrpcRequests.Inc()

	return &productsService.CreateProductRes{ProductID: result.ProductID.String()}, nil
}

func (s *ProductGrpcServiceServer) UpdateProduct(ctx context.Context, req *productsService.UpdateProductReq) (*productsService.UpdateProductRes, error) {
	s.Metrics.UpdateProductGrpcRequests.Inc()
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "ProductGrpcServiceServer.UpdateProductCommand")
	span.LogFields(log.String("req", req.String()))
	defer span.Finish()

	productUUID, err := uuid.FromString(req.GetProductID())
	if err != nil {
		s.Log.WarnMsg("uuid.FromString", err)
		tracing.TraceErr(span, err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	command := updatingProductV1.NewUpdateProductCommand(productUUID, req.GetName(), req.GetDescription(), req.GetPrice())

	if err := s.Validator.StructCtx(ctx, command); err != nil {
		s.Log.WarnMsg("validate", err)
		tracing.TraceErr(span, err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	if _, err = mediatr.Send[*updatingProductV1.UpdateProductCommand, *mediatr.Unit](ctx, command); err != nil {
		s.Log.WarnMsg("UpdateProductCommand.Handle", err)
		tracing.TraceErr(span, err)
		return nil, s.errResponse(codes.Internal, err)
	}

	s.Log.Infof("(product updated) id: {%s}", productUUID.String())
	s.Metrics.SuccessGrpcRequests.Inc()

	return &productsService.UpdateProductRes{}, nil
}

func (s *ProductGrpcServiceServer) GetProductById(ctx context.Context, req *productsService.GetProductByIdReq) (*productsService.GetProductByIdRes, error) {

	s.Metrics.GetProductByIdGrpcRequests.Inc()
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "ProductGrpcServiceServer.GetProductByIdQuery")
	span.LogFields(log.String("req", req.String()))
	defer span.Finish()

	productUUID, err := uuid.FromString(req.GetProductID())
	if err != nil {
		s.Log.WarnMsg("uuid.FromString", err)
		tracing.TraceErr(span, err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	query := &gettingProductByIdV1.GetProductByIdQuery{ProductID: productUUID}
	if err := s.Validator.StructCtx(ctx, query); err != nil {
		s.Log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	queryResult, err := mediatr.Send[*gettingProductByIdV1.GetProductByIdQuery, *gettingProductByIdDtos.GetProductByIdResponseDto](ctx, query)
	if err != nil {
		s.Log.WarnMsg("GetProductByIdQuery.Handle", err)
		tracing.TraceErr(span, err)
		return nil, s.errResponse(codes.Internal, err)
	}

	product, err := mapper.Map[*models.Product](queryResult.Product)

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, s.errResponse(codes.Internal, err)
	}

	s.Metrics.SuccessGrpcRequests.Inc()

	pr, err := mapper.Map[*productsService.Product](product)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, s.errResponse(codes.Internal, err)
	}

	return &productsService.GetProductByIdRes{Product: pr}, nil
}

func (s *ProductGrpcServiceServer) errResponse(c codes.Code, err error) error {
	s.Metrics.ErrorGrpcRequests.Inc()
	return status.Error(c, err.Error())
}
