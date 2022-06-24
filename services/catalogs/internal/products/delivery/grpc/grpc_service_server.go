package grpc

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	product_service_client "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/contracts/grpc/service_clients"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/contracts/repositories"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/features/creating_product"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/features/getting_product_by_id"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/features/updating_product"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/mappers"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/models"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/shared/configurations"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductGrpcServiceServer struct {
	mediator          *mediatr.Mediator
	infrastructure    *configurations.Infrastructure
	productRepository repositories.ProductRepository
	// Ref:https://github.com/grpc/grpc-go/issues/3794#issuecomment-720599532
	// product_service_client.UnimplementedProductsServiceServer
}

func NewProductGrpcService(infra *configurations.Infrastructure, mediator *mediatr.Mediator, productRepository repositories.ProductRepository) *ProductGrpcServiceServer {
	return &ProductGrpcServiceServer{productRepository: productRepository, infrastructure: infra, mediator: mediator}
}

func (s *ProductGrpcServiceServer) CreateProduct(ctx context.Context, req *product_service_client.CreateProductReq) (*product_service_client.CreateProductRes, error) {
	s.infrastructure.Metrics.CreateProductGrpcRequests.Inc()

	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "grpcService.CreateProduct")
	defer span.Finish()

	command := creating_product.NewCreateProduct(req.GetName(), req.GetDescription(), req.GetPrice())
	if err := s.infrastructure.Validator.StructCtx(ctx, command); err != nil {
		s.infrastructure.Log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	_, err := s.mediator.Send(ctx, command)
	if err != nil {
		s.infrastructure.Log.WarnMsg("CreateProduct.Handle", err)
		return nil, s.errResponse(codes.Internal, err)
	}

	s.infrastructure.Metrics.SuccessGrpcRequests.Inc()
	return &product_service_client.CreateProductRes{ProductID: command.ProductID.String()}, nil
}

func (s *ProductGrpcServiceServer) UpdateProduct(ctx context.Context, req *product_service_client.UpdateProductReq) (*product_service_client.UpdateProductRes, error) {
	s.infrastructure.Metrics.UpdateProductGrpcRequests.Inc()

	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "grpcService.UpdateProduct")
	defer span.Finish()

	productUUID, err := uuid.FromString(req.GetProductID())
	if err != nil {
		s.infrastructure.Log.WarnMsg("uuid.FromString", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	command := updating_product.NewUpdateProduct(productUUID, req.GetName(), req.GetDescription(), req.GetPrice())
	if err := s.infrastructure.Validator.StructCtx(ctx, command); err != nil {
		s.infrastructure.Log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	_, err = s.mediator.Send(ctx, command)
	if err != nil {
		s.infrastructure.Log.WarnMsg("UpdateProduct.Handle", err)
		return nil, s.errResponse(codes.Internal, err)
	}

	s.infrastructure.Metrics.SuccessGrpcRequests.Inc()
	return &product_service_client.UpdateProductRes{}, nil
}

func (s *ProductGrpcServiceServer) GetProductById(ctx context.Context, req *product_service_client.GetProductByIdReq) (*product_service_client.GetProductByIdRes, error) {
	s.infrastructure.Metrics.GetProductByIdGrpcRequests.Inc()

	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "grpcService.GetProductById")
	defer span.Finish()

	productUUID, err := uuid.FromString(req.GetProductID())
	if err != nil {
		s.infrastructure.Log.WarnMsg("uuid.FromString", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	query := getting_product_by_id.NewGetProductById(productUUID)
	if err := s.infrastructure.Validator.StructCtx(ctx, query); err != nil {
		s.infrastructure.Log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	queryResult, err := s.mediator.Send(ctx, query)
	if err != nil {
		s.infrastructure.Log.WarnMsg("GetProductById.Handle", err)
		return nil, s.errResponse(codes.Internal, err)
	}

	p, ok := queryResult.(models.Product)
	if err := utils.CheckType(ok); err != nil {
		return nil, err
	}

	s.infrastructure.Metrics.SuccessGrpcRequests.Inc()
	return &product_service_client.GetProductByIdRes{Product: mappers.WriterProductToGrpc(&p)}, nil
}

func (s *ProductGrpcServiceServer) errResponse(c codes.Code, err error) error {
	s.infrastructure.Metrics.ErrorGrpcRequests.Inc()
	return status.Error(c, err.Error())
}
