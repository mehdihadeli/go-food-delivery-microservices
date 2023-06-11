package fxlog

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
)

var FxLogger = fx.WithLogger(func(logger logger.Logger) fxevent.Logger {
	return NewCustomFxLogger(logger)
},
)

// Ref: https://articles.wesionary.team/logging-interfaces-in-go-182c28be3d18

type FxCustomLogger struct {
	logger.Logger
}

func NewCustomFxLogger(logger logger.Logger) fxevent.Logger {
	return &FxCustomLogger{Logger: logger}
}

// Printf prits go-fxlog logs
func (l FxCustomLogger) Printf(str string, args ...interface{}) {
	if len(args) > 0 {
		l.Debugf(str, args)
	}
	l.Debug(str)
}

func (l *FxCustomLogger) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		l.Debug("OnStart hook executing: ",
			zap.String("callee", e.FunctionName),
			zap.String("caller", e.CallerName),
		)
	case *fxevent.OnStartExecuted:
		if e.Err != nil {
			l.Debug("OnStart hook failed: ",
				zap.String("callee", e.FunctionName),
				zap.String("caller", e.CallerName),
				zap.Error(e.Err),
			)
		} else {
			l.Debug("OnStart hook executed: ",
				zap.String("callee", e.FunctionName),
				zap.String("caller", e.CallerName),
				zap.String("runtime", e.Runtime.String()),
			)
		}
	case *fxevent.OnStopExecuting:
		l.Debug("OnStop hook executing: ",
			zap.String("callee", e.FunctionName),
			zap.String("caller", e.CallerName),
		)
	case *fxevent.OnStopExecuted:
		if e.Err != nil {
			l.Debug("OnStop hook failed: ",
				zap.String("callee", e.FunctionName),
				zap.String("caller", e.CallerName),
				zap.Error(e.Err),
			)
		} else {
			l.Debug("OnStop hook executed: ",
				zap.String("callee", e.FunctionName),
				zap.String("caller", e.CallerName),
				zap.String("runtime", e.Runtime.String()),
			)
		}
	case *fxevent.Supplied:
		l.Debug("supplied: ", zap.String("type", e.TypeName), zap.Error(e.Err))
	case *fxevent.Provided:
		for _, rtype := range e.OutputTypeNames {
			l.Debug("provided: ", e.ConstructorName, " => ", rtype)
		}
	case *fxevent.Decorated:
		for _, rtype := range e.OutputTypeNames {
			l.Debug("decorated: ",
				zap.String("decorator", e.DecoratorName),
				zap.String("type", rtype),
			)
		}
	case *fxevent.Invoking:
		l.Debug("invoking: ", e.FunctionName)
	case *fxevent.Started:
		if e.Err == nil {
			l.Debug("started")
		}
	case *fxevent.LoggerInitialized:
		if e.Err == nil {
			l.Debug("initialized: custom fxevent.Logger -> ", e.ConstructorName)
		}
	}
}
