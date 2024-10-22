package middleware

import (
	"go.uber.org/zap"
)

// ZapLogger адаптер для zap логгера
type ZapLogger struct {
	logger *zap.Logger
}

func NewZapLogger() (*ZapLogger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	return &ZapLogger{
		logger: logger,
	}, nil
}

func (l *ZapLogger) Info(args ...interface{}) {
	l.logger.Sugar().Info(args...)
}

func (l *ZapLogger) Error(args ...interface{}) {
	l.logger.Sugar().Error(args...)
}

func (l *ZapLogger) Infof(template string, args ...interface{}) {
	l.logger.Sugar().Infof(template, args...)
}

func (l *ZapLogger) Errorf(template string, args ...interface{}) {
	l.logger.Sugar().Errorf(template, args...)
}
