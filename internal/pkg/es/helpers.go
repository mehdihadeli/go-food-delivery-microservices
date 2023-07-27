package es

import (
	"fmt"

	"go.uber.org/fx"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/es/contracts/projection"
)

func AsProjection(handler interface{}) interface{} {
	return fx.Annotate(
		handler,
		fx.As(new(projection.IProjection)),
		fx.ResultTags(fmt.Sprintf(`group:"projections"`)),
	)
}
