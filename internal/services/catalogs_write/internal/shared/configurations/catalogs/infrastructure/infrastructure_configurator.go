package infrastructure

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/fxapp"
)

type InfrastructureConfigurator struct {
	*fxapp.Application
}

func NewInfrastructureConfigurator(fxapp *fxapp.Application) *InfrastructureConfigurator {
	return &InfrastructureConfigurator{
		Application: fxapp,
	}
}

func (ic *InfrastructureConfigurator) ConfigInfrastructures() {
	ic.ResolveFunc(func() error {
		return nil
	})
}
