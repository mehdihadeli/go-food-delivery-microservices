package grpc

import (
	"context"
	"github.com/adzeitor/mediatr"
	"github.com/go-playground/validator"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/features/creating_product"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/features/getting_product_by_id"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/features/updating_product"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/mappers"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/models"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/proto/product_service"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/shared"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductGrpcService struct {
	log     logger.Logger
	cfg     *config.Config
	v       *validator.Validate
	md      *mediatr.Mediator
	metrics *shared.CatalogsServiceMetrics
	// Ref:https://github.com/grpc/grpc-go/issues/3794#issuecomment-720599532
	product_service.UnimplementedProductsServiceServer
}

func NewProductGrpcService(log logger.Logger, cfg *config.Config, v *validator.Validate, md *mediatr.Mediator, metrics *shared.CatalogsServiceMetrics) *ProductGrpcService {
	return &ProductGrpcService{log: log, cfg: cfg, v: v, md: md, metrics: metrics}
}

func (s *ProductGrpcService) CreateProduct(ctx context.Context, req *product_service.CreateProductReq) (*product_service.CreateProductRes, error) {
	s.metrics.CreateProductGrpcRequests.Inc()

	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "grpcService.CreateProduct")
	defer span.Finish()

	productUUID, err := uuid.FromString(req.GetProductID())
	if err != nil {
		s.log.WarnMsg("uuid.FromString", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	command := creating_product.NewCreateProduct(productUUID, req.GetName(), req.GetDescription(), req.GetPrice())
	if err := s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	_, err = s.md.Send(ctx, command)
	if err != nil {
		s.log.WarnMsg("CreateProduct.Handle", err)
		return nil, s.errResponse(codes.Internal, err)
	}

	s.metrics.SuccessGrpcRequests.Inc()
	return &product_service.CreateProductRes{ProductID: productUUID.String()}, nil
}

func (s *ProductGrpcService) UpdateProduct(ctx context.Context, req *product_service.UpdateProductReq) (*product_service.UpdateProductRes, error) {
	s.metrics.UpdateProductGrpcRequests.Inc()

	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "grpcService.UpdateProduct")
	defer span.Finish()

	productUUID, err := uuid.FromString(req.GetProductID())
	if err != nil {
		s.log.WarnMsg("uuid.FromString", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	command := updating_product.NewUpdateProduct(productUUID, req.GetName(), req.GetDescription(), req.GetPrice())
	if err := s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	_, err = s.md.Send(ctx, command)
	if err != nil {
		s.log.WarnMsg("UpdateProduct.Handle", err)
		return nil, s.errResponse(codes.Internal, err)
	}

	s.metrics.SuccessGrpcRequests.Inc()
	return &product_service.UpdateProductRes{}, nil
}

func (s *ProductGrpcService) GetProductById(ctx context.Context, req *product_service.GetProductByIdReq) (*product_service.GetProductByIdRes, error) {
	s.metrics.GetProductByIdGrpcRequests.Inc()

	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "grpcService.GetProductById")
	defer span.Finish()

	productUUID, err := uuid.FromString(req.GetProductID())
	if err != nil {
		s.log.WarnMsg("uuid.FromString", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	query := getting_product_by_id.NewGetProductById(productUUID)
	if err := s.v.StructCtx(ctx, query); err != nil {
		s.log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	queryResult, err := s.md.Send(ctx, query)
	if err != nil {
		s.log.WarnMsg("GetProductById.Handle", err)
		return nil, s.errResponse(codes.Internal, err)
	}

	p, ok := queryResult.(models.Product)
	if err := errors.CheckType(ok); err != nil {
		return nil, err
	}

	s.metrics.SuccessGrpcRequests.Inc()
	return &product_service.GetProductByIdRes{Product: mappers.WriterProductToGrpc(&p)}, nil
}

func (s *ProductGrpcService) errResponse(c codes.Code, err error) error {
	s.metrics.ErrorGrpcRequests.Inc()
	return status.Error(c, err.Error())
}
