package http

import (
	"go.uber.org/fx"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/client"
	customEcho "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/custom_echo"
)

// Module provided to fxlog
// https://uber-go.github.io/fx/modules.html
var Module = fx.Module("httpfx",
	// - order is not important in provide
	// - provide can have parameter and will resolve if registered
	// - execute its func only if it requested
	client.Module,
	customEcho.Module,
)
