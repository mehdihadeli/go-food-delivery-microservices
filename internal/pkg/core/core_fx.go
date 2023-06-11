package core

import (
	"go.uber.org/fx"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/serializer"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/serializer/json"
)

// Module provided to fxlog
// https://uber-go.github.io/fx/modules.html
var Module = fx.Module(
	"corefx",
	fx.Provide(
		json.NewJsonSerializer,
		serializer.NewDefaultEventSerializer,
		serializer.NewDefaultMetadataSerializer,
	),
)
