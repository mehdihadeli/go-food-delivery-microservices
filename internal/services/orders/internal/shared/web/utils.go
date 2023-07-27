package web

import (
    "fmt"
    "strings"

    "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orders/config"
)

func GetMicroserviceName(cfg *config.Config) string {
	return fmt.Sprintf("(%s)", strings.ToUpper(cfg.ServiceName))
}
