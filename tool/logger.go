package tool

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
)

var logger zapLogger
var loggerOnce sync.Once

type zapLogger struct {
	logger *zap.Logger
	sugar  *zap.SugaredLogger
}

func Logger() zapLogger {
	if logger.logger == nil {
		loggerOnce.Do(func() {
			logger = zapLogger{}
			config := zap.NewProductionConfig()
			encoderConfig := zap.NewProductionEncoderConfig()
			encoderConfig.TimeKey = "timestamp"
			encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
			config.EncoderConfig = encoderConfig
			newLogger, _ := config.Build(zap.AddCallerSkip(1))
			logger.logger = newLogger
			newSugar := newLogger.Sugar()
			logger.sugar = newSugar
		})
	}
	return logger
}

func (l zapLogger) Info(msg string, params ...string) {
	defer l.logger.Sync()
	l.sugar.Infow(msg, zap.Strings("params", params))
}

func (l zapLogger) Warning(msg string, err error, params ...string) {
	defer l.logger.Sync()
	l.sugar.Warnw(
		msg,
		zap.Error(err),
		zap.Strings("params", params),
	)
}

func (l zapLogger) Error(msg string, err error, params ...string) {
	defer l.logger.Sync()
	l.sugar.Error(
		msg,
		zap.Error(err),
		zap.Strings("params", params),
	)
}
