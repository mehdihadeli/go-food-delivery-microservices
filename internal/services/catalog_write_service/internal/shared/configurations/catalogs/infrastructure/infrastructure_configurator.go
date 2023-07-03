package infrastructure

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/fxapp/contracts"
)

type InfrastructureConfigurator struct {
	contracts.Application
}

func NewInfrastructureConfigurator(fxapp contracts.Application) *InfrastructureConfigurator {
	return &InfrastructureConfigurator{
		Application: fxapp,
	}
}

func (ic *InfrastructureConfigurator) ConfigInfrastructures() {
	ic.ResolveFunc(func() error {
		return nil
	})
}
