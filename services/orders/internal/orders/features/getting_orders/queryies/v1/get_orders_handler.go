package v1

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/config"
)

// TODO: Should read from read side model (mongo)

type GetOrdersQueryHandler struct {
	log logger.Logger
	cfg *config.Config
}
