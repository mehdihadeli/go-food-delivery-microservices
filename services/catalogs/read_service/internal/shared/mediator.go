package shared

import (
	kafkaClient "github.com/mehdihadeli/store-golang-microservice-sample/pkg/kafka"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/config"
	"github.com/pkg/errors"
)

func NewCatalogsMediator(log logger.Logger, cfg *config.Config, kafkaProducer kafkaClient.Producer) (*mediatr.Mediator, error) {

	md := mediatr.New()

	err := md.Register()

	if err != nil {
		return nil, errors.Wrap(err, "error while registering handlers in the mediator")
	}

	return &md, nil
}
