package fxlog

import (
	"strings"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

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

// LogEvent logs the given event to the provided Zap logger.
func (l *FxCustomLogger) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		l.Debugw("OnStart hook executing", logger.Fields{"caller": e.CallerName, "function": e.FunctionName})
	case *fxevent.OnStartExecuted:
		if e.Err != nil {
			l.Errorw("OnStart hook failed",
				logger.Fields{"caller": e.CallerName, "callee": e.CallerName, "error": e.Err},
			)
		} else {
			l.Debugw("OnStart hook executed", logger.Fields{"caller": e.CallerName, "callee": e.FunctionName, "runtime": e.Runtime.String()})
		}
	case *fxevent.OnStopExecuting:
		l.Debugw("OnStop hook executing", logger.Fields{"callee": e.FunctionName, "caller": e.CallerName})
	case *fxevent.OnStopExecuted:
		if e.Err != nil {
			l.Errorw("OnStop hook failed",
				logger.Fields{"caller": e.CallerName, "callee": e.CallerName, "error": e.Err},
			)
		} else {
			l.Debugw("OnStop hook executed", logger.Fields{"caller": e.CallerName, "callee": e.FunctionName, "runtime": e.Runtime.String()})
		}
	case *fxevent.Supplied:
		if e.Err != nil {
			l.Errorw("error encountered while applying options",
				logger.Fields{"type": e.TypeName, "stacktrace": e.StackTrace, "module": e.ModuleName, "error": e.Err},
			)
		} else {
			l.Debugw("supplied", logger.Fields{"type": e.TypeName, "stacktrace": e.StackTrace, "module": e.ModuleName})
		}
	case *fxevent.Provided:
		for _, rtype := range e.OutputTypeNames {
			l.Debugw("provided", logger.Fields{"constructor": e.ConstructorName, "stacktrace": e.StackTrace, "module": e.ModuleName, "type": rtype, "private": e.Private})
		}
		if e.Err != nil {
			l.Errorw("error encountered while applying options",
				logger.Fields{"module": e.ModuleName, "stacktrace": e.StackTrace, "error": e.Err},
			)
		}
	case *fxevent.Replaced:
		for _, rtype := range e.OutputTypeNames {
			l.Debugw("replaced", logger.Fields{"stacktrace": e.StackTrace, "module": e.ModuleName, "type": rtype})
		}
		if e.Err != nil {
			l.Errorw("error encountered while replacing",
				logger.Fields{"module": e.ModuleName, "stacktrace": e.StackTrace, "error": e.Err},
			)
		}
	case *fxevent.Decorated:
		for _, rtype := range e.OutputTypeNames {
			l.Debugw("decorated", logger.Fields{"decorator": e.DecoratorName, "stacktrace": e.StackTrace, "module": e.ModuleName, "type": rtype})
		}
		if e.Err != nil {
			l.Errorw("error encountered while applying options",
				logger.Fields{"module": e.ModuleName, "stacktrace": e.StackTrace, "error": e.Err},
			)
		}
	case *fxevent.Run:
		if e.Err != nil {
			l.Errorw("error returned",
				logger.Fields{"module": e.ModuleName, "name": e.Name, "kind": e.Kind, "error": e.Err},
			)
		} else {
			l.Debugw("run", logger.Fields{"module": e.ModuleName, "name": e.Name, "kind": e.Kind})
		}
	case *fxevent.Invoking:
		// Do not log stack as it will make logs hard to read.
		l.Debugw("invoking", logger.Fields{"module": e.ModuleName, "function": e.FunctionName})
	case *fxevent.Invoked:
		if e.Err != nil {
			l.Errorw("invoke failed",
				logger.Fields{"error": e.Err, "stack": e.Trace, "function": e.FunctionName, "module": e.ModuleName},
			)
		}
	case *fxevent.Stopping:
		l.Debugw("received signal", logger.Fields{"signal": strings.ToUpper(e.Signal.String())})
	case *fxevent.Stopped:
		if e.Err != nil {
			l.Errorw("stop failed",
				logger.Fields{"error": e.Err},
			)
		}
	case *fxevent.RollingBack:
		l.Errorw("start failed, rolling back",
			logger.Fields{"error": e.StartErr},
		)
	case *fxevent.RolledBack:
		if e.Err != nil {
			l.Errorw("rollback failed",
				logger.Fields{"error": e.Err},
			)
		}
	case *fxevent.Started:
		if e.Err != nil {
			l.Errorw("start failed",
				logger.Fields{"error": e.Err},
			)
		} else {
			l.Debug("started")
		}
	case *fxevent.LoggerInitialized:
		if e.Err != nil {
			l.Errorw("custom logger initialization failed",
				logger.Fields{"error": e.Err},
			)
		} else {
			l.Debugw("initialized custom fxevent.Logger", logger.Fields{"function": e.ConstructorName})
		}
	}
}
