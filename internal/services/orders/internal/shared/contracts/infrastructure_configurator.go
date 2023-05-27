package contracts

import (
	"context"
)

type InfrastructureConfigurator interface {
	ConfigInfrastructures(ctx context.Context) (*InfrastructureConfigurations, func(), error)
}
