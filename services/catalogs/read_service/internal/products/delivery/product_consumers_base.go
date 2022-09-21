package delivery

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/configurations/infrastructure"
)

type ProductConsumersBase struct {
	*infrastructure.InfrastructureConfigurations
}

func NewProductConsumersBase(infra *infrastructure.InfrastructureConfigurations) *ProductConsumersBase {
	return &ProductConsumersBase{InfrastructureConfigurations: infra}
}

func (pm *ProductConsumersBase) CommitMessage() {
	pm.Metrics.SuccessKafkaMessages.Inc()
}

func (pm *ProductConsumersBase) CommitErrMessage() {
	pm.Metrics.ErrorKafkaMessages.Inc()
}
