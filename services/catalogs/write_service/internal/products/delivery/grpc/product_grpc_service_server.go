package grpc

import (
	"context"
	"fmt"

	"emperror.dev/errors"
	"github.com/mehdihadeli/go-mediatr"
	uuid "github.com/satori/go.uuid"
	attribute2 "go.opentelemetry.io/otel/attribute"

	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mapper"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/bus"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/otel/tracing/attribute"
	productsService "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts/proto/service_clients"
	createProductCommandV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/creating_product/v1/commands"
	createProductDtosV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/creating_product/v1/dtos"
	getProductByIdDtosV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_product_by_id/v1/dtos"
	getProductByIdQueryV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_product_by_id/v1/queries"
	updateProductCommandV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/updating_product/v1/commands"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/contracts"
	"go.opentelemetry.io/otel/trace"
)

var grpcMetricsAttr = attribute2.Key("MetricsType").String("Grpc")

type ProductGrpcServiceServer struct {
	*contracts.InfrastructureConfigurations
	catalogsMetrics *contracts.CatalogsMetrics
	bus             bus.Bus
	// Ref:https://github.com/grpc/grpc-go/issues/3794#issuecomment-720599532
	// product_service_client.UnimplementedProductsServiceServer
}

func NewProductGrpcService(infra *contracts.InfrastructureConfigurations, catalogsMetrics *contracts.CatalogsMetrics, bus bus.Bus) *ProductGrpcServiceServer {
	return &ProductGrpcServiceServer{InfrastructureConfigurations: infra, catalogsMetrics: catalogsMetrics, bus: bus}
}

func (s *ProductGrpcServiceServer) CreateProduct(ctx context.Context, req *productsService.CreateProductReq) (*productsService.CreateProductRes, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.Object("Request", req))
	s.catalogsMetrics.CreateProductGrpcRequests.Add(ctx, 1, grpcMetricsAttr)

	command, err := createProductCommandV1.NewCreateProduct(req.GetName(), req.GetDescription(), req.GetPrice())
	if err != nil {
		validationErr := customErrors.NewValidationErrorWrap(err, "[ProductGrpcServiceServer_CreateProduct.StructCtx] command validation failed")
		s.Log.Errorf(fmt.Sprintf("[ProductGrpcServiceServer_CreateProduct.StructCtx] err: %v", validationErr))
		return nil, validationErr
	}

	result, err := mediatr.Send[*createProductCommandV1.CreateProduct, *createProductDtosV1.CreateProductResponseDto](ctx, command)
	if err != nil {
		err = errors.WithMessage(err, "[ProductGrpcServiceServer_CreateProduct.Send] error in sending CreateProduct")
		s.Log.Errorw(fmt.Sprintf("[ProductGrpcServiceServer_CreateProduct.Send] id: {%s}, err: %v", command.ProductID, err), logger.Fields{"ProductId": command.ProductID})
		return nil, err
	}

	return &productsService.CreateProductRes{ProductId: result.ProductID.String()}, nil
}

func (s *ProductGrpcServiceServer) UpdateProduct(ctx context.Context, req *productsService.UpdateProductReq) (*productsService.UpdateProductRes, error) {
	s.catalogsMetrics.UpdateProductGrpcRequests.Add(ctx, 1, grpcMetricsAttr)
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.Object("Request", req))

	productUUID, err := uuid.FromString(req.GetProductId())
	if err != nil {
		badRequestErr := customErrors.NewBadRequestErrorWrap(err, "[ProductGrpcServiceServer_UpdateProduct.uuid.FromString] error in converting uuid")
		s.Log.Errorf(fmt.Sprintf("[ProductGrpcServiceServer_UpdateProduct.uuid.FromString] err: %v", badRequestErr))
		return nil, badRequestErr
	}

	command, err := updateProductCommandV1.NewUpdateProduct(productUUID, req.GetName(), req.GetDescription(), req.GetPrice())
	if err != nil {
		validationErr := customErrors.NewValidationErrorWrap(err, "[ProductGrpcServiceServer_UpdateProduct.StructCtx] command validation failed")
		s.Log.Errorf(fmt.Sprintf("[ProductGrpcServiceServer_UpdateProduct.StructCtx] err: %v", validationErr))
		return nil, validationErr
	}

	if _, err = mediatr.Send[*updateProductCommandV1.UpdateProduct, *mediatr.Unit](ctx, command); err != nil {
		err = errors.WithMessage(err, "[ProductGrpcServiceServer_UpdateProduct.Send] error in sending CreateProduct")
		s.Log.Errorw(fmt.Sprintf("[ProductGrpcServiceServer_UpdateProduct.Send] id: {%s}, err: %v", command.ProductID, err), logger.Fields{"ProductId": command.ProductID})
		return nil, err
	}

	return &productsService.UpdateProductRes{}, nil
}

func (s *ProductGrpcServiceServer) GetProductById(ctx context.Context, req *productsService.GetProductByIdReq) (*productsService.GetProductByIdRes, error) {
	//// we could use trace manually, but I used grpc middleware for doing this
	//ctx, span, clean := grpcTracing.StartGrpcServerTracerSpan(ctx, "ProductGrpcServiceServer.GetProductById")
	//defer clean()

	s.catalogsMetrics.GetProductByIdGrpcRequests.Add(ctx, 1, grpcMetricsAttr)
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.Object("Request", req))

	productUUID, err := uuid.FromString(req.GetProductId())
	if err != nil {
		badRequestErr := customErrors.NewBadRequestErrorWrap(err, "[ProductGrpcServiceServer_GetProductById.uuid.FromString] error in converting uuid")
		s.Log.Errorf(fmt.Sprintf("[ProductGrpcServiceServer_GetProductById.uuid.FromString] err: %v", badRequestErr))
		return nil, badRequestErr
	}

	query, err := getProductByIdQueryV1.NewGetProductById(productUUID)
	if err != nil {
		validationErr := customErrors.NewValidationErrorWrap(err, "[ProductGrpcServiceServer_GetProductById.StructCtx] query validation failed")
		s.Log.Errorf(fmt.Sprintf("[ProductGrpcServiceServer_GetProductById.StructCtx] err: %v", validationErr))
		return nil, validationErr
	}

	queryResult, err := mediatr.Send[*getProductByIdQueryV1.GetProductById, *getProductByIdDtosV1.GetProductByIdResponseDto](ctx, query)
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

	return &productsService.GetProductByIdRes{Product: product}, nil
}
