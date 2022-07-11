package grpc

import (
	"context"

	getting_product_by_id_dtos "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_product_by_id/dtos"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/models"

	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mapper"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mediatr"

	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	product_service_client "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts/grpc/service_clients"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/creating_product"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/creating_product/dtos"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_product_by_id"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/updating_product"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/mappings"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/configurations/infrastructure"
	"github.com/opentracing/opentracing-go/log"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductGrpcServiceServer struct {
	infrastructure *infrastructure.InfrastructureConfiguration
	// Ref:https://github.com/grpc/grpc-go/issues/3794#issuecomment-720599532
	// product_service_client.UnimplementedProductsServiceServer
}

func NewProductGrpcService(infra *infrastructure.InfrastructureConfiguration) *ProductGrpcServiceServer {
	return &ProductGrpcServiceServer{infrastructure: infra}
}

func (s *ProductGrpcServiceServer) CreateProduct(ctx context.Context, req *product_service_client.CreateProductReq) (*product_service_client.CreateProductRes, error) {

	s.infrastructure.Metrics.CreateProductGrpcRequests.Inc()
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "productGrpcServiceServer.CreateProduct")
	span.LogFields(log.String("req", req.String()))
	defer span.Finish()

	command := creating_product.NewCreateProduct(req.GetName(), req.GetDescription(), req.GetPrice())

	if err := s.infrastructure.Validator.StructCtx(ctx, command); err != nil {
		s.infrastructure.Log.Errorf("(validate) err: {%v}", err)
		tracing.TraceErr(span, err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	result, err := mediatr.Send[*dtos.CreateProductResponseDto](ctx, command)
	if err != nil {
		s.infrastructure.Log.Errorf("(CreateProduct.Handle) productId: {%s}, err: {%v}", command.ProductID, err)
		tracing.TraceErr(span, err)
		return nil, s.errResponse(codes.Internal, err)
	}

	s.infrastructure.Log.Infof("(product created) productId: {%s}", command.ProductID)
	s.infrastructure.Metrics.SuccessGrpcRequests.Inc()

	return &product_service_client.CreateProductRes{ProductID: result.ProductID.String()}, nil
}

func (s *ProductGrpcServiceServer) UpdateProduct(ctx context.Context, req *product_service_client.UpdateProductReq) (*product_service_client.UpdateProductRes, error) {
	s.infrastructure.Metrics.UpdateProductGrpcRequests.Inc()
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "productGrpcServiceServer.UpdateProduct")
	span.LogFields(log.String("req", req.String()))
	defer span.Finish()

	productUUID, err := uuid.FromString(req.GetProductID())
	if err != nil {
		s.infrastructure.Log.WarnMsg("uuid.FromString", err)
		tracing.TraceErr(span, err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	command := updating_product.NewUpdateProduct(productUUID, req.GetName(), req.GetDescription(), req.GetPrice())

	if err := s.infrastructure.Validator.StructCtx(ctx, command); err != nil {
		s.infrastructure.Log.WarnMsg("validate", err)
		tracing.TraceErr(span, err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	_, err = mediatr.Send[*mediatr.Unit](ctx, command)
	if err != nil {
		s.infrastructure.Log.WarnMsg("UpdateProduct.Handle", err)
		tracing.TraceErr(span, err)
		return nil, s.errResponse(codes.Internal, err)
	}

	s.infrastructure.Log.Infof("(product updated) id: {%s}", productUUID.String())
	s.infrastructure.Metrics.SuccessGrpcRequests.Inc()

	return &product_service_client.UpdateProductRes{}, nil
}

func (s *ProductGrpcServiceServer) GetProductById(ctx context.Context, req *product_service_client.GetProductByIdReq) (*product_service_client.GetProductByIdRes, error) {

	s.infrastructure.Metrics.GetProductByIdGrpcRequests.Inc()
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "grpcService.GetProductById")
	span.LogFields(log.String("req", req.String()))
	defer span.Finish()

	productUUID, err := uuid.FromString(req.GetProductID())
	if err != nil {
		s.infrastructure.Log.WarnMsg("uuid.FromString", err)
		tracing.TraceErr(span, err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	query := getting_product_by_id.NewGetProductById(productUUID)
	if err := s.infrastructure.Validator.StructCtx(ctx, query); err != nil {
		s.infrastructure.Log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	queryResult, err := mediatr.Send[*getting_product_by_id_dtos.GetProductByIdResponseDto](ctx, query)
	if err != nil {
		s.infrastructure.Log.WarnMsg("GetProductById.Handle", err)
		tracing.TraceErr(span, err)
		return nil, s.errResponse(codes.Internal, err)
	}

	product, err := mapper.Map[*models.Product](queryResult.Product)

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, s.errResponse(codes.Internal, err)
	}

	s.infrastructure.Metrics.SuccessGrpcRequests.Inc()

	return &product_service_client.GetProductByIdRes{Product: mappings.WriterProductToGrpc(product)}, nil
}

func (s *ProductGrpcServiceServer) errResponse(c codes.Code, err error) error {
	s.infrastructure.Metrics.ErrorGrpcRequests.Inc()
	return status.Error(c, err.Error())
}
