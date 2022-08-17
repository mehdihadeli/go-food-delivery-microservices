package zap

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"os"
	"time"

	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/constants"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	level       string
	encoding    string
	sugarLogger *zap.SugaredLogger
	logger      *zap.Logger
}

type ZapLogger interface {
	logger.Logger
	DPanic(args ...interface{})
	DPanicf(template string, args ...interface{})
	Sync() error
}

// For mapping config logger
var loggerLevelMap = map[string]zapcore.Level{
	"debug": zapcore.DebugLevel,
	"info":  zapcore.InfoLevel,
	"warn":  zapcore.WarnLevel,
	"error": zapcore.ErrorLevel,
	"panic": zapcore.PanicLevel,
	"fatal": zapcore.FatalLevel,
}
var (
	DefaultLogger logger.Logger
)

// NewZapLogger create new zap logger
func NewZapLogger(cfg *logger.Config) ZapLogger {
	zapLogger := &zapLogger{level: cfg.LogLevel, encoding: cfg.Encoder}
	zapLogger.initLogger()

	return zapLogger
}

func init() {
	DefaultLogger = NewZapLogger(&logger.Config{
		LogLevel: "debug",
		Encoder:  "json",
	})
}

func (l *zapLogger) getLoggerLevel() zapcore.Level {
	level, exist := loggerLevelMap[l.level]
	if !exist {
		return zapcore.DebugLevel
	}

	return level
}

// InitLogger Init logger
func (l *zapLogger) initLogger() {
	logLevel := l.getLoggerLevel()

	logWriter := zapcore.AddSync(os.Stdout)

	var encoderCfg zapcore.EncoderConfig

	env := os.Getenv("APP_ENV")
	if env == constants.Production {
		encoderCfg = zap.NewProductionEncoderConfig()
	} else {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	}

	var encoder zapcore.Encoder
	encoderCfg.NameKey = "[SERVICE]"
	encoderCfg.TimeKey = "[TIME]"
	encoderCfg.LevelKey = "[LEVEL]"
	encoderCfg.FunctionKey = "[CALLER]"
	encoderCfg.CallerKey = "[LINE]"
	encoderCfg.MessageKey = "[MESSAGE]"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder
	encoderCfg.EncodeName = zapcore.FullNameEncoder
	encoderCfg.EncodeDuration = zapcore.StringDurationEncoder

	if l.encoding == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}

	core := zapcore.NewCore(encoder, logWriter, zap.NewAtomicLevelAt(logLevel))
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	l.logger = logger
	l.sugarLogger = logger.Sugar()
}

// Logger methods

// WithName add logger microservice name
func (l *zapLogger) WithName(name string) {
	l.logger = l.logger.Named(name)
	l.sugarLogger = l.sugarLogger.Named(name)
}

// Debug uses fmt.Sprint to construct and log a message.
func (l *zapLogger) Debug(args ...interface{}) {
	l.sugarLogger.Debug(args...)
}

// Debugf uses fmt.Sprintf to log a templated message
func (l *zapLogger) Debugf(template string, args ...interface{}) {
	l.sugarLogger.Debugf(template, args...)
}

// Info uses fmt.Sprint to construct and log a message
func (l *zapLogger) Info(args ...interface{}) {
	l.sugarLogger.Info(args...)
}

// Infof uses fmt.Sprintf to log a templated message.
func (l *zapLogger) Infof(template string, args ...interface{}) {
	l.sugarLogger.Infof(template, args...)
}

// Printf uses fmt.Sprintf to log a templated message
func (l *zapLogger) Printf(template string, args ...interface{}) {
	l.sugarLogger.Infof(template, args...)
}

// Warn uses fmt.Sprint to construct and log a message.
func (l *zapLogger) Warn(args ...interface{}) {
	l.sugarLogger.Warn(args...)
}

// WarnMsg log error message with warn level.
func (l *zapLogger) WarnMsg(msg string, err error) {
	l.logger.Warn(msg, zap.String("error", err.Error()))
}

// Warnf uses fmt.Sprintf to log a templated message.
func (l *zapLogger) Warnf(template string, args ...interface{}) {
	l.sugarLogger.Warnf(template, args...)
}

// Error uses fmt.Sprint to construct and log a message.
func (l *zapLogger) Error(args ...interface{}) {
	l.sugarLogger.Error(args...)
}

// Errorf uses fmt.Sprintf to log a templated message.
func (l *zapLogger) Errorf(template string, args ...interface{}) {
	l.sugarLogger.Errorf(template, args...)
}

// Err uses error to log a message.
func (l *zapLogger) Err(msg string, err error) {
	l.logger.Error(msg, zap.Error(err))
}

// DPanic uses fmt.Sprint to construct and log a message. In development, the logger then panics. (See DPanicLevel for details.)
func (l *zapLogger) DPanic(args ...interface{}) {
	l.sugarLogger.DPanic(args...)
}

// DPanicf uses fmt.Sprintf to log a templated message. In development, the logger then panics. (See DPanicLevel for details.)
func (l *zapLogger) DPanicf(template string, args ...interface{}) {
	l.sugarLogger.DPanicf(template, args...)
}

// Panic uses fmt.Sprint to construct and log a message, then panics.
func (l *zapLogger) Panic(args ...interface{}) {
	l.sugarLogger.Panic(args...)
}

// Panicf uses fmt.Sprintf to log a templated message, then panics
func (l *zapLogger) Panicf(template string, args ...interface{}) {
	l.sugarLogger.Panicf(template, args...)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func (l *zapLogger) Fatal(args ...interface{}) {
	l.sugarLogger.Fatal(args...)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func (l *zapLogger) Fatalf(template string, args ...interface{}) {
	l.sugarLogger.Fatalf(template, args...)
}

// Sync flushes any buffered log entries
func (l *zapLogger) Sync() error {
	go func() {
		err := l.logger.Sync()
		if err != nil {
			l.logger.Error("error while syncing", zap.Error(err))
		}
	}() // nolint: errcheck
	return l.sugarLogger.Sync()
}

func (l *zapLogger) HttpMiddlewareAccessLogger(method, uri string, status int, size int64, time time.Duration) {
	l.Info(
		constants.HTTP,
		zap.String(constants.METHOD, method),
		zap.String(constants.URI, uri),
		zap.Int(constants.STATUS, status),
		zap.Int64(constants.SIZE, size),
		zap.Duration(constants.TIME, time),
	)
}

func (l *zapLogger) GrpcMiddlewareAccessLogger(method string, time time.Duration, metaData map[string][]string, err error) {
	l.Info(
		constants.GRPC,
		zap.String(constants.METHOD, method),
		zap.Duration(constants.TIME, time),
		zap.Any(constants.METADATA, metaData),
		zap.Error(err),
	)
}

func (l *zapLogger) GrpcClientInterceptorLogger(method string, req, reply interface{}, time time.Duration, metaData map[string][]string, err error) {
	l.Info(
		constants.GRPC,
		zap.String(constants.METHOD, method),
		zap.Any(constants.REQUEST, req),
		zap.Any(constants.REPLY, reply),
		zap.Duration(constants.TIME, time),
		zap.Any(constants.METADATA, metaData),
		zap.Error(err),
	)
}

func (l *zapLogger) KafkaProcessMessage(topic string, partition int, message string, workerID int, offset int64, time time.Time) {
	l.Debug(
		"Processing Kafka message",
		zap.String(constants.Topic, topic),
		zap.Int(constants.Partition, partition),
		zap.String(constants.Message, message),
		zap.Int(constants.WorkerID, workerID),
		zap.Int64(constants.Offset, offset),
		zap.Time(constants.Time, time),
	)
}

func (l *zapLogger) KafkaLogCommittedMessage(topic string, partition int, offset int64) {
	l.Info(
		"Committed Kafka message",
		zap.String(constants.Topic, topic),
		zap.Int(constants.Partition, partition),
		zap.Int64(constants.Offset, offset),
	)
}
