package contracts

import "context"

type InfrastructureConfigurator interface {
	ConfigInfrastructures(ctx context.Context) (InfrastructureConfiguration, error, func())
}
