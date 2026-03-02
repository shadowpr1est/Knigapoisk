package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger() *zap.Logger {
	cfg := zap.NewProductionConfig()
	cfg.Encoding = "json"
	cfg.EncoderConfig.TimeKey = "ts"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stderr"}

	l, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(l)
	return l
}

func WithService(l *zap.Logger, service string) *zap.Logger {
	return l.With(zap.String("service", service), zap.String("env", os.Getenv("APP_ENV")))
}

