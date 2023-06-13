package infrastructure

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/fxapp"
)

type InfrastructureConfigurator struct {
	// log       logger.Logger
	// cfg       *config.AppOptions
	// validator *validator.Validate
	// meter     metric.Meter
	*fxapp.Application
}

func NewInfrastructureConfigurator(fxapp *fxapp.Application,

// log logger.Logger,
// cfg *config.AppOptions,
// validator *validator.Validate,
// meter metric.Meter,
) *InfrastructureConfigurator {
	return &InfrastructureConfigurator{
		// log:       log,
		// cfg:       cfg,
		// validator: validator,
		// meter:     meter,
		Application: fxapp,
	}
}

func (ic *InfrastructureConfigurator) ConfigInfrastructures() {
	ic.ResolveFunc(func() error {
		return nil
	})
}
