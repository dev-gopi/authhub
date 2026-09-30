package logger

import (
	"fmt"

	"go.uber.org/zap"
)

var Log *zap.Logger

func InitLogger(environment string) error {
	logger, err := New(environment)
	if err != nil {
		return fmt.Errorf("initialize logger: %w", err)
	}

	Log = logger

	Log.Info(
		"logger initialized",
		zap.String("environment", environment),
	)

	return nil
}

func New(environment string) (*zap.Logger, error) {
	if environment == "production" {
		logger, err := zap.NewProduction()
		if err != nil {
			return nil, fmt.Errorf("create production logger: %w", err)
		}

		return logger, nil
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, fmt.Errorf("create development logger: %w", err)
	}

	return logger, nil
}

func Info(msg string, fields ...zap.Field) {
	if Log == nil {
		return
	}

	Log.Info(msg, fields...)
}

func Error(msg string, err error, fields ...zap.Field) {
	if Log == nil {
		return
	}

	if err != nil {
		fields = append(
			fields,
			zap.Error(err),
		)
	}

	Log.Error(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	if Log == nil {
		return
	}

	Log.Warn(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	if Log == nil {
		return
	}

	Log.Debug(msg, fields...)
}

func Sync() {
	if Log == nil {
		return
	}

	_ = Log.Sync()
}
