package contracts

import (
	"go.opentelemetry.io/otel/metric"
)

type CatalogsMetrics struct {
	CreateProductGrpcRequests     metric.Float64Counter
	UpdateProductGrpcRequests     metric.Float64Counter
	DeleteProductGrpcRequests     metric.Float64Counter
	GetProductByIdGrpcRequests    metric.Float64Counter
	SearchProductGrpcRequests     metric.Float64Counter
	SuccessRabbitMQMessages       metric.Float64Counter
	ErrorRabbitMQMessages         metric.Float64Counter
	CreateProductRabbitMQMessages metric.Float64Counter
	UpdateProductRabbitMQMessages metric.Float64Counter
	DeleteProductRabbitMQMessages metric.Float64Counter
}
