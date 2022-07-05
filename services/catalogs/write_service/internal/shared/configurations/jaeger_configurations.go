package configurations

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/opentracing/opentracing-go"
)

func (ic *infrastructureConfigurator) configJaeger() (error, func()) {
	if ic.cfg.Jaeger.Enable {
		tracer, closer, err := tracing.NewJaegerTracer(ic.cfg.Jaeger)
		if err != nil {
			return err, nil
		}
		opentracing.SetGlobalTracer(tracer)
		return nil, func() {
			_ = closer.Close()
		}
	}

	return nil, func() {}
}
