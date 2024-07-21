package empty

import (
	"time"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger/models"
)

var EmptyLogger logger.Logger = &emptyLogger{}

type emptyLogger struct{}

func (e emptyLogger) Configure(cfg func(internalLog interface{})) {
}

func (e emptyLogger) Debug(args ...interface{}) {
}

func (e emptyLogger) Debugf(template string, args ...interface{}) {
}

func (e emptyLogger) Debugw(msg string, fields logger.Fields) {
}

func (e emptyLogger) LogType() models.LogType {
	return models.Zap
}

func (e emptyLogger) Info(args ...interface{}) {
}

func (e emptyLogger) Infof(template string, args ...interface{}) {
}

func (e emptyLogger) Infow(msg string, fields logger.Fields) {
}

func (e emptyLogger) Warn(args ...interface{}) {
}

func (e emptyLogger) Warnf(template string, args ...interface{}) {
}

func (e emptyLogger) WarnMsg(msg string, err error) {
}

func (e emptyLogger) Error(args ...interface{}) {
}

func (e emptyLogger) Errorw(msg string, fields logger.Fields) {
}

func (e emptyLogger) Errorf(template string, args ...interface{}) {
}

func (e emptyLogger) Err(msg string, err error) {
}

func (e emptyLogger) Fatal(args ...interface{}) {
}

func (e emptyLogger) Fatalf(template string, args ...interface{}) {
}

func (e emptyLogger) Printf(template string, args ...interface{}) {
}

func (e emptyLogger) WithName(name string) {
}

func (e emptyLogger) GrpcMiddlewareAccessLogger(
	method string,
	time time.Duration,
	metaData map[string][]string,
	err error,
) {
}

func (e emptyLogger) GrpcClientInterceptorLogger(
	method string,
	req interface{},
	reply interface{},
	time time.Duration,
	metaData map[string][]string,
	err error,
) {
}
