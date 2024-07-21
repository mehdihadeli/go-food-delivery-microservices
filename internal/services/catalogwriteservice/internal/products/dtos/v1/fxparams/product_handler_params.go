package fxparams

import (
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/messaging/producer"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/otel/tracing"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/shared/data/dbcontext"

	"go.uber.org/fx"
)

type ProductHandlerParams struct {
	fx.In

	Log               logger.Logger
	CatalogsDBContext *dbcontext.CatalogsGormDBContext
	RabbitmqProducer  producer.Producer
	Tracer            tracing.AppTracer
}
