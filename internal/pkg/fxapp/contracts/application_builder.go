package contracts

import (
	"go.uber.org/fx"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config/environemnt"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
)

type ApplicationBuilder interface {
	ProvideModule(module fx.Option)
	Provide(constructors ...interface{})
	Decorate(constructors ...interface{})
	Build() Application

	GetProvides() []interface{}
	GetDecorates() []interface{}
	Options() []fx.Option
	Logger() logger.Logger
	Environment() environemnt.Environment
}
