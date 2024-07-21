package loggingpipelines

import (
	"context"
	"fmt"
	"time"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger"
	typeMapper "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/reflection/typemapper"

	"github.com/mehdihadeli/go-mediatr"
)

type requestLoggerPipeline struct {
	logger logger.Logger
}

func NewMediatorLoggingPipeline(l logger.Logger) mediatr.PipelineBehavior {
	return &requestLoggerPipeline{logger: l}
}

func (r *requestLoggerPipeline) Handle(
	ctx context.Context,
	request interface{},
	next mediatr.RequestHandlerFunc,
) (interface{}, error) {
	startTime := time.Now()
	defer func() {
		elapsed := time.Since(startTime)
		r.logger.Infof("Request took %s", elapsed)
	}()

	requestName := typeMapper.GetNonePointerTypeName(request)

	r.logger.Infow(
		fmt.Sprintf("Handling request: '%s'", requestName),
		logger.Fields{"Request": request},
	)

	response, err := next(ctx)
	if err != nil {
		r.logger.Infof("Request failed with error: %v", err)

		return nil, err
	}

	responseName := typeMapper.GetNonePointerTypeName(response)

	r.logger.Infow(
		fmt.Sprintf(
			"Request handled successfully with response: '%s'",
			responseName,
		),
		logger.Fields{"Response": response},
	)

	return response, nil
}
