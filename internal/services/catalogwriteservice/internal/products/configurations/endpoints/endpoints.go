package endpoints

import (
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/web/route"
)

func RegisterEndpoints(endpoints []route.Endpoint) error {
	for _, endpoint := range endpoints {
		endpoint.MapEndpoint()
	}

	return nil
}
