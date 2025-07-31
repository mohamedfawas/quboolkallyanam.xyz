package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Init(debug bool) (*zap.Logger, error) {
	var cfg zap.Config
	if debug {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	logger, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}
	return logger, nil
}
