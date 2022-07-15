package web

import (
	"fmt"
	"strings"
	"thub.com/mehdihadeli/store-golang-microservice-sample/services/orders/config"
)

func GetMicroserviceName(cfg *config.Config) string {
	return fmt.Sprintf("(%s)", strings.ToUpper(cfg.ServiceName))
}
