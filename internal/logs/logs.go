package logs

import (
	"log"

	"go.uber.org/zap"
)

var Logger *zap.Logger

func InitLogger(develMode bool) {
	var err error

	if develMode {
		Logger, err = zap.NewDevelopment()
	} else {
		cfg := zap.NewProductionConfig()
		cfg.DisableCaller = true
		cfg.DisableStacktrace = true
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		Logger, err = cfg.Build()
	}
	if err != nil {
		log.Fatal("Error: cannot init zap", err)
	}
}

func Error(text string, fields ...zap.Field) {
	Logger.Error(text, fields...)
}

func Info(text string, fields ...zap.Field) {
	Logger.Info(text, fields...)
}
