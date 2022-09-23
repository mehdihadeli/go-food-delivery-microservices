package consumers

import "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/configurations/infrastructure"

func ConfigConsumers(infra *infrastructure.InfrastructureConfiguration) error {
	//add custom message type mappings
	//utils.RegisterCustomMessageTypesToRegistrty(map[string]types.IMessage{"orderCreatedV1": &OrderCreatedV1{}})

	return nil
}
