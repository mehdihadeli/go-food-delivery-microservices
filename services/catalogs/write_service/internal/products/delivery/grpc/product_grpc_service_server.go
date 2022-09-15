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
	productsService "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts/proto/service_clients"
	creatingProductV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/creating_product/commands/v1"
	gettingProductByIdV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_product_by_id/queries/v1"
	updatingProductV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/updating_product/commands/v1"

	creatingProductDtos "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/creating_product/dtos"
	gettingProductByIdDtos "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_product_by_id/dtos"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/configurations/infrastructure"
	"github.com/opentracing/opentracing-go/log"
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
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "ProductGrpcServiceServer.CreateProduct")
	span.LogFields(log.Object("Request", req))
	s.Metrics.CreateProductGrpcRequests.Inc()
	defer span.Finish()

	command := creatingProductV1.NewCreateProductCommand(req.GetName(), req.GetDescription(), req.GetPrice())

	if err := s.Validator.StructCtx(ctx, command); err != nil {
		validationErr := customErrors.NewValidationErrorWrap(err, "[ProductGrpcServiceServer_CreateProduct.StructCtx] command validation failed")
		s.Log.Errorf(fmt.Sprintf("[ProductGrpcServiceServer_CreateProduct.StructCtx] err: %v", tracing.TraceWithErr(span, validationErr)))
		return nil, grpcErrors.ErrGrpcResponse(validationErr)
	}

	result, err := mediatr.Send[*creatingProductV1.CreateProductCommand, *creatingProductDtos.CreateProductResponseDto](ctx, command)
	if err != nil {
		err = errors.WithMessage(err, "[ProductGrpcServiceServer_CreateProduct.Send] error in sending CreateProductCommand")
		s.Log.Errorw(fmt.Sprintf("[ProductGrpcServiceServer_CreateProduct.Send] id: {%s}, err: %v", command.ProductID, tracing.TraceWithErr(span, err)), logger.Fields{"ProductId": command.ProductID})
		return nil, grpcErrors.ErrGrpcResponse(err)
	}

	s.Metrics.SuccessGrpcRequests.Inc()

	return &productsService.CreateProductRes{ProductID: result.ProductID.String()}, nil
}

func (s *ProductGrpcServiceServer) UpdateProduct(ctx context.Context, req *productsService.UpdateProductReq) (*productsService.UpdateProductRes, error) {
	s.Metrics.UpdateProductGrpcRequests.Inc()
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "ProductGrpcServiceServer.UpdateProduct")
	span.LogFields(log.Object("Request", req))
	defer span.Finish()

	productUUID, err := uuid.FromString(req.GetProductID())
	if err != nil {
		badRequestErr := customErrors.NewBadRequestErrorWrap(err, "[ProductGrpcServiceServer_UpdateProduct.uuid.FromString] error in converting uuid")
		s.Log.Errorf(fmt.Sprintf("[ProductGrpcServiceServer_UpdateProduct.uuid.FromString] err: %v", tracing.TraceWithErr(span, badRequestErr)))
		return nil, grpcErrors.ErrGrpcResponse(badRequestErr)
	}

	command := updatingProductV1.NewUpdateProductCommand(productUUID, req.GetName(), req.GetDescription(), req.GetPrice())

	if err := s.Validator.StructCtx(ctx, command); err != nil {
		validationErr := customErrors.NewValidationErrorWrap(err, "[ProductGrpcServiceServer_UpdateProduct.StructCtx] command validation failed")
		s.Log.Errorf(fmt.Sprintf("[ProductGrpcServiceServer_UpdateProduct.StructCtx] err: %v", tracing.TraceWithErr(span, validationErr)))
		return nil, grpcErrors.ErrGrpcResponse(validationErr)
	}

	if _, err = mediatr.Send[*updatingProductV1.UpdateProductCommand, *mediatr.Unit](ctx, command); err != nil {
		err = errors.WithMessage(err, "[ProductGrpcServiceServer_UpdateProduct.Send] error in sending CreateProductCommand")
		s.Log.Errorw(fmt.Sprintf("[ProductGrpcServiceServer_UpdateProduct.Send] id: {%s}, err: %v", command.ProductID, tracing.TraceWithErr(span, err)), logger.Fields{"ProductId": command.ProductID})
		return nil, grpcErrors.ErrGrpcResponse(err)
	}

	s.Metrics.SuccessGrpcRequests.Inc()

	return &productsService.UpdateProductRes{}, nil
}

func (s *ProductGrpcServiceServer) GetProductById(ctx context.Context, req *productsService.GetProductByIdReq) (*productsService.GetProductByIdRes, error) {
	s.Metrics.GetProductByIdGrpcRequests.Inc()
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "ProductGrpcServiceServer.GetProductById")
	span.LogFields(log.Object("Request", req))
	defer span.Finish()

	productUUID, err := uuid.FromString(req.GetProductID())
	if err != nil {
		badRequestErr := customErrors.NewBadRequestErrorWrap(err, "[ProductGrpcServiceServer_GetProductById.uuid.FromString] error in converting uuid")
		s.Log.Errorf(fmt.Sprintf("[ProductGrpcServiceServer_GetProductById.uuid.FromString] err: %v", tracing.TraceWithErr(span, badRequestErr)))
		return nil, grpcErrors.ErrGrpcResponse(badRequestErr)
	}

	query := gettingProductByIdV1.NewGetProductByIdQuery(productUUID)
	if err := s.Validator.StructCtx(ctx, query); err != nil {
		validationErr := customErrors.NewValidationErrorWrap(err, "[ProductGrpcServiceServer_GetProductById.StructCtx] query validation failed")
		s.Log.Errorf(fmt.Sprintf("[ProductGrpcServiceServer_GetProductById.StructCtx] err: %v", tracing.TraceWithErr(span, validationErr)))
		return nil, grpcErrors.ErrGrpcResponse(validationErr)
	}

	queryResult, err := mediatr.Send[*gettingProductByIdV1.GetProductByIdQuery, *gettingProductByIdDtos.GetProductByIdResponseDto](ctx, query)
	if err != nil {
		err = errors.WithMessage(err, "[ProductGrpcServiceServer_GetProductById.Send] error in sending GetProductByIdQuery")
		s.Log.Errorw(fmt.Sprintf("[ProductGrpcServiceServer_GetProductById.Send] id: {%s}, err: %v", query.ProductID, tracing.TraceWithErr(span, err)), logger.Fields{"ProductId": query.ProductID})
		return nil, grpcErrors.ErrGrpcResponse(err)
	}

	product, err := mapper.Map[*productsService.Product](queryResult.Product)
	if err != nil {
		err = errors.WithMessage(err, "[ProductGrpcServiceServer_GetProductById.Map] error in mapping product")
		return nil, grpcErrors.ErrGrpcResponse(tracing.TraceWithErr(span, err))
	}

	s.Metrics.SuccessGrpcRequests.Inc()

	return &productsService.GetProductByIdRes{Product: product}, nil
}
