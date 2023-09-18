package gromlog

import (
	"context"
	"time"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"

	gormlogger "gorm.io/gorm/logger"
)

// Ref: https://articles.wesionary.team/logging-interfaces-in-go-182c28be3d18
// implement gorm logger Interface

type GormCustomLogger struct {
	logger.Logger
	gormlogger.Config
}

func NewGormCustomLogger(logger logger.Logger) *GormCustomLogger {
	//cfg, err := config.ProvideLogConfig()
	//
	//var logger logger.Logger
	//if cfg.LogType == models.Logrus && err != nil {
	//	logger = logrous.NewLogrusLogger(cfg, constants.Dev)
	//} else {
	//	if err != nil {
	//		cfg = &config.LogOptions{LogLevel: "info", LogType: models.Zap}
	//	}
	//	logger = zap.NewZapLogger(cfg, constants.Dev)
	//}

	return &GormCustomLogger{
		Logger: logger,
		Config: gormlogger.Config{
			LogLevel: gormlogger.Info,
		},
	}
}

// LogMode set log mode
func (l *GormCustomLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

// Info prints info
func (l GormCustomLogger) Info(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel >= gormlogger.Info {
		l.Debugf(str, args...)
	}
}

// Warn prints warn messages
func (l GormCustomLogger) Warn(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel >= gormlogger.Warn {
		l.Warnf(str, args...)
	}
}

// Error prints error messages
func (l GormCustomLogger) Error(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel >= gormlogger.Error {
		l.Errorf(str, args...)
	}
}

// Trace prints trace messages
func (l GormCustomLogger) Trace(
	ctx context.Context,
	begin time.Time,
	fc func() (string, int64),
	err error,
) {
	if l.LogLevel <= 0 {
		return
	}
	elapsed := time.Since(begin)
	if l.LogLevel >= gormlogger.Info {
		sql, rows := fc()
		l.Debug("[", elapsed.Milliseconds(), " ms, ", rows, " rows] ", "sql -> ", sql)
		return
	}

	if l.LogLevel >= gormlogger.Warn {
		sql, rows := fc()
		l.Logger.Warn("[", elapsed.Milliseconds(), " ms, ", rows, " rows] ", "sql -> ", sql)
		return
	}

	if l.LogLevel >= gormlogger.Error {
		sql, rows := fc()
		l.Logger.Error("[", elapsed.Milliseconds(), " ms, ", rows, " rows] ", "sql -> ", sql)
		return
	}
}
