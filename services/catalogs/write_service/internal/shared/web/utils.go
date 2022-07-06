package web

import (
	"fmt"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/config"
	"strings"
)

func GetMicroserviceName(cfg *config.Config) string {
	return fmt.Sprintf("(%s)", strings.ToUpper(cfg.ServiceName))
}
