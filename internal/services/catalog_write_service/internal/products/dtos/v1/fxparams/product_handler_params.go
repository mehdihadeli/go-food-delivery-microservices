package fxparams

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/producer"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/contracts"

	"go.uber.org/fx"
)

type ProductHandlerParams struct {
	fx.In

	Log               logger.Logger
	Uow               contracts.CatalogUnitOfWork
	ProductRepository contracts.ProductRepository
	RabbitmqProducer  producer.Producer
	Tracer            tracing.AppTracer
}
