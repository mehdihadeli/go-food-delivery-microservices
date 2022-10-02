package grpc

import (
	"context"
	"emperror.dev/errors"
	"fmt"
	"github.com/mehdihadeli/go-mediatr"
	grpcTracing "github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc/otel/tracing"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mapper"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/otel/tracing/attribute"
	productsService "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts/proto/service_clients"
	creatingProductV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/creating_product/commands/v1"
	gettingProductByIdDtos "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_product_by_id/dtos"
	gettingProductByIdV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_product_by_id/queries/v1"
	updatingProductV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/updating_product/commands/v1"

	creatingProductDtos "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/creating_product/dtos"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/configurations/infrastructure"
	uuid "github.com/satori/go.uuid"
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
	span := grpcTracing.SpanFromContext(ctx)
	span.SetAttributes(attribute.Object("Request", req))
	s.Metrics.CreateProductGrpcRequests.Inc()

	command := creatingProductV1.NewCreateProduct(req.GetName(), req.GetDescription(), req.GetPrice())

	if err := s.Validator.StructCtx(ctx, command); err != nil {
		validationErr := customErrors.NewValidationErrorWrap(err, "[ProductGrpcServiceServer_CreateProduct.StructCtx] command validation failed")
		s.Log.Errorf(fmt.Sprintf("[ProductGrpcServiceServer_CreateProduct.StructCtx] err: %v", validationErr))
		return nil, validationErr
	}

	result, err := mediatr.Send[*creatingProductV1.CreateProduct, *creatingProductDtos.CreateProductResponseDto](ctx, command)
	if err != nil {
		err = errors.WithMessage(err, "[ProductGrpcServiceServer_CreateProduct.Send] error in sending CreateProduct")
		s.Log.Errorw(fmt.Sprintf("[ProductGrpcServiceServer_CreateProduct.Send] id: {%s}, err: %v", command.ProductID, err), logger.Fields{"ProductId": command.ProductID})
		return nil, err
	}

	s.Metrics.SuccessGrpcRequests.Inc()

	return &productsService.CreateProductRes{ProductId: result.ProductID.String()}, nil
}

func (s *ProductGrpcServiceServer) UpdateProduct(ctx context.Context, req *productsService.UpdateProductReq) (*productsService.UpdateProductRes, error) {
	s.Metrics.UpdateProductGrpcRequests.Inc()
	span := grpcTracing.SpanFromContext(ctx)
	span.SetAttributes(attribute.Object("Request", req))

	productUUID, err := uuid.FromString(req.GetProductId())
	if err != nil {
		badRequestErr := customErrors.NewBadRequestErrorWrap(err, "[ProductGrpcServiceServer_UpdateProduct.uuid.FromString] error in converting uuid")
		s.Log.Errorf(fmt.Sprintf("[ProductGrpcServiceServer_UpdateProduct.uuid.FromString] err: %v", badRequestErr))
		return nil, badRequestErr
	}

	command := updatingProductV1.NewUpdateProduct(productUUID, req.GetName(), req.GetDescription(), req.GetPrice())

	if err := s.Validator.StructCtx(ctx, command); err != nil {
		validationErr := customErrors.NewValidationErrorWrap(err, "[ProductGrpcServiceServer_UpdateProduct.StructCtx] command validation failed")
		s.Log.Errorf(fmt.Sprintf("[ProductGrpcServiceServer_UpdateProduct.StructCtx] err: %v", validationErr))
		return nil, validationErr
	}

	if _, err = mediatr.Send[*updatingProductV1.UpdateProduct, *mediatr.Unit](ctx, command); err != nil {
		err = errors.WithMessage(err, "[ProductGrpcServiceServer_UpdateProduct.Send] error in sending CreateProduct")
		s.Log.Errorw(fmt.Sprintf("[ProductGrpcServiceServer_UpdateProduct.Send] id: {%s}, err: %v", command.ProductID, err), logger.Fields{"ProductId": command.ProductID})
		return nil, err
	}

	s.Metrics.SuccessGrpcRequests.Inc()

	return &productsService.UpdateProductRes{}, nil
}

func (s *ProductGrpcServiceServer) GetProductById(ctx context.Context, req *productsService.GetProductByIdReq) (*productsService.GetProductByIdRes, error) {
	//// we could use trace manually, but I used grpc middleware for doing this
	//ctx, span, clean := grpcTracing.StartGrpcServerTracerSpan(ctx, "ProductGrpcServiceServer.GetProductById")
	//defer clean()

	s.Metrics.GetProductByIdGrpcRequests.Inc()
	span := grpcTracing.SpanFromContext(ctx)
	span.SetAttributes(attribute.Object("Request", req))

	productUUID, err := uuid.FromString(req.GetProductId())
	if err != nil {
		badRequestErr := customErrors.NewBadRequestErrorWrap(err, "[ProductGrpcServiceServer_GetProductById.uuid.FromString] error in converting uuid")
		s.Log.Errorf(fmt.Sprintf("[ProductGrpcServiceServer_GetProductById.uuid.FromString] err: %v", badRequestErr))
		return nil, badRequestErr
	}

	query := gettingProductByIdV1.NewGetProductById(productUUID)
	if err := s.Validator.StructCtx(ctx, query); err != nil {
		validationErr := customErrors.NewValidationErrorWrap(err, "[ProductGrpcServiceServer_GetProductById.StructCtx] query validation failed")
		s.Log.Errorf(fmt.Sprintf("[ProductGrpcServiceServer_GetProductById.StructCtx] err: %v", validationErr))
		return nil, validationErr
	}

	queryResult, err := mediatr.Send[*gettingProductByIdV1.GetProductById, *gettingProductByIdDtos.GetProductByIdResponseDto](ctx, query)
	if err != nil {
		err = errors.WithMessage(err, "[ProductGrpcServiceServer_GetProductById.Send] error in sending GetProductById")
		s.Log.Errorw(fmt.Sprintf("[ProductGrpcServiceServer_GetProductById.Send] id: {%s}, err: %v", query.ProductID, err), logger.Fields{"ProductId": query.ProductID})
		return nil, err
	}

	product, err := mapper.Map[*productsService.Product](queryResult.Product)
	if err != nil {
		err = errors.WithMessage(err, "[ProductGrpcServiceServer_GetProductById.Map] error in mapping product")
		return nil, err
	}

	s.Metrics.SuccessGrpcRequests.Inc()

	return &productsService.GetProductByIdRes{Product: product}, nil
}
