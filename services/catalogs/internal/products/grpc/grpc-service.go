package grpc

import (
	"context"
	"github.com/go-playground/validator"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/config"
	product_service "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/grpc/proto/product-service"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/mappers"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/metrics"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/service"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcService struct {
	log     logger.Logger
	cfg     *config.Config
	v       *validator.Validate
	ps      *service.ProductService
	metrics *metrics.ProductServiceMetrics
}

func NewWriterGrpcService(log logger.Logger, cfg *config.Config, v *validator.Validate, ps *service.ProductService, metrics *metrics.ProductServiceMetrics) *grpcService {
	return &grpcService{log: log, cfg: cfg, v: v, ps: ps, metrics: metrics}
}

func (s *grpcService) CreateProduct(ctx context.Context, req *product_service.CreateProductReq) (*product_service.CreateProductRes, error) {
	s.metrics.CreateProductGrpcRequests.Inc()

	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "grpcService.CreateProduct")
	defer span.Finish()

	productUUID, err := uuid.FromString(req.GetProductID())
	if err != nil {
		s.log.WarnMsg("uuid.FromString", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	command := commands.NewCreateProductCommand(productUUID, req.GetName(), req.GetDescription(), req.GetPrice())
	if err := s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	err = s.ps.Commands.CreateProduct.Handle(ctx, command)
	if err != nil {
		s.log.WarnMsg("CreateProduct.Handle", err)
		return nil, s.errResponse(codes.Internal, err)
	}

	s.metrics.SuccessGrpcRequests.Inc()
	return &product_service.CreateProductRes{ProductID: productUUID.String()}, nil
}

func (s *grpcService) UpdateProduct(ctx context.Context, req *product_service.UpdateProductReq) (*product_service.UpdateProductRes, error) {
	s.metrics.UpdateProductGrpcRequests.Inc()

	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "grpcService.UpdateProduct")
	defer span.Finish()

	productUUID, err := uuid.FromString(req.GetProductID())
	if err != nil {
		s.log.WarnMsg("uuid.FromString", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	command := commands.NewUpdateProductCommand(productUUID, req.GetName(), req.GetDescription(), req.GetPrice())
	if err := s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	err = s.ps.Commands.UpdateProduct.Handle(ctx, command)
	if err != nil {
		s.log.WarnMsg("UpdateProduct.Handle", err)
		return nil, s.errResponse(codes.Internal, err)
	}

	s.metrics.SuccessGrpcRequests.Inc()
	return &product_service.UpdateProductRes{}, nil
}

func (s *grpcService) GetProductById(ctx context.Context, req *product_service.GetProductByIdReq) (*product_service.GetProductByIdRes, error) {
	s.metrics.GetProductByIdGrpcRequests.Inc()

	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "grpcService.GetProductById")
	defer span.Finish()

	productUUID, err := uuid.FromString(req.GetProductID())
	if err != nil {
		s.log.WarnMsg("uuid.FromString", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	query := queries.NewGetProductByIdQuery(productUUID)
	if err := s.v.StructCtx(ctx, query); err != nil {
		s.log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	product, err := s.ps.Queries.GetProductById.Handle(ctx, query)
	if err != nil {
		s.log.WarnMsg("GetProductById.Handle", err)
		return nil, s.errResponse(codes.Internal, err)
	}

	s.metrics.SuccessGrpcRequests.Inc()
	return &product_service.GetProductByIdRes{Product: mappers.WriterProductToGrpc(product)}, nil
}

func (s *grpcService) errResponse(c codes.Code, err error) error {
	s.metrics.ErrorGrpcRequests.Inc()
	return status.Error(c, err.Error())
}
