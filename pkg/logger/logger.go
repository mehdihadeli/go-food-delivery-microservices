package logger

import (
	"time"
)

type Config struct {
	LogLevel string `mapstructure:"level"`
	Encoder  string `mapstructure:"encoder"`
}

func NewLoggerConfig(logLevel string, encoder string) *Config {
	return &Config{LogLevel: logLevel, Encoder: encoder}
}

type Logger interface {
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	WarnMsg(msg string, err error)
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Err(msg string, err error)
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
	Printf(template string, args ...interface{})
	WithName(name string)
	HttpMiddlewareAccessLogger(method string, uri string, status int, size int64, time time.Duration)
	GrpcMiddlewareAccessLogger(method string, time time.Duration, metaData map[string][]string, err error)
	GrpcClientInterceptorLogger(method string, req interface{}, reply interface{}, time time.Duration, metaData map[string][]string, err error)
	KafkaProcessMessage(topic string, partition int, message string, workerID int, offset int64, time time.Time)
	KafkaLogCommittedMessage(topic string, partition int, offset int64)
}
