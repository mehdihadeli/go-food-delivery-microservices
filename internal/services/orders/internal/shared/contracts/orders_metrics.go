package contracts

import (
	"go.opentelemetry.io/otel/metric"
)

type OrdersMetrics struct {
	SuccessGrpcRequests metric.Float64Counter
	ErrorGrpcRequests   metric.Float64Counter

	CreateOrderGrpcRequests  metric.Float64Counter
	UpdateOrderGrpcRequests  metric.Float64Counter
	PayOrderGrpcRequests     metric.Float64Counter
	SubmitOrderGrpcRequests  metric.Float64Counter
	GetOrderByIdGrpcRequests metric.Float64Counter
	GetOrdersGrpcRequests    metric.Float64Counter
	SearchOrderGrpcRequests  metric.Float64Counter

	SuccessHttpRequests metric.Float64Counter
	ErrorHttpRequests   metric.Float64Counter

	CreateOrderHttpRequests  metric.Float64Counter
	UpdateOrderHttpRequests  metric.Float64Counter
	PayOrderHttpRequests     metric.Float64Counter
	SubmitOrderHttpRequests  metric.Float64Counter
	GetOrderByIdHttpRequests metric.Float64Counter
	SearchOrderHttpRequests  metric.Float64Counter
	GetOrdersHttpRequests    metric.Float64Counter

	SuccessRabbitMQMessages metric.Float64Counter
	ErrorRabbitMQMessages   metric.Float64Counter

	CreateOrderRabbitMQMessages metric.Float64Counter
	UpdateOrderRabbitMQMessages metric.Float64Counter
	DeleteOrderRabbitMQMessages metric.Float64Counter
}
