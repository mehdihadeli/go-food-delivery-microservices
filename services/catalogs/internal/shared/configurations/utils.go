package configurations

import (
	"fmt"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/config"
	"strings"
)

func GetMicroserviceName(cfg *config.Config) string {
	return fmt.Sprintf("(%s)", strings.ToUpper(cfg.ServiceName))
}
