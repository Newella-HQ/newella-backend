package logger

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/Newella-HQ/newella-backend/internal/config"
)

type Logger interface {
	With(args ...interface{}) Logger

	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})

	Debugln(args ...interface{})
	Infoln(args ...interface{})
	Warnln(args ...interface{})
	Errorln(args ...interface{})
	Fatalln(args ...interface{})

	Sync() error
}

type ZapLogger struct {
	logger *zap.SugaredLogger
}

func NewZapLogger(level config.LogLevel) (*ZapLogger, error) {
	cfg := zap.NewProductionConfig()

	switch level {
	case config.Debug:
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case config.Info:
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case config.Warn:
		cfg.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case config.Error:
		cfg.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	}

	l, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("can't init loggger: %w", err)
	}

	return &ZapLogger{
		logger: l.Sugar().With("log_level", l.Level()),
	}, nil
}

func (l *ZapLogger) With(args ...interface{}) Logger {
	return &ZapLogger{
		logger: l.logger.With(args),
	}
}

func (l *ZapLogger) Debugf(template string, args ...interface{}) {
	l.logger.Debugf(template, args...)
}

func (l *ZapLogger) Infof(template string, args ...interface{}) {
	l.logger.Infof(template, args...)
}

func (l *ZapLogger) Warnf(template string, args ...interface{}) {
	l.logger.Warnf(template, args...)
}

func (l *ZapLogger) Errorf(template string, args ...interface{}) {
	l.logger.Errorf(template, args...)
}

func (l *ZapLogger) Fatalf(template string, args ...interface{}) {
	l.logger.Fatalf(template, args...)
}

func (l *ZapLogger) Debugln(args ...interface{}) {
	l.logger.Debugln(args)
}

func (l *ZapLogger) Infoln(args ...interface{}) {
	l.logger.Infoln(args)
}

func (l *ZapLogger) Warnln(args ...interface{}) {
	l.logger.Warnln(args)
}

func (l *ZapLogger) Errorln(args ...interface{}) {
	l.logger.Errorln(args)
}

func (l *ZapLogger) Fatalln(args ...interface{}) {
	l.logger.Fatalln(args)
}

func (l *ZapLogger) Sync() error {
	return l.logger.Sync()
}
