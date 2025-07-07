package logger

import (
	"os"

	"go.uber.org/zap"
)

var Log *zap.SugaredLogger

func InitLogger() {
	env := os.Getenv("ENVIRONMENT")

	var zapLogger *zap.Logger
	var err error

	if env == "production" {
		zapLogger, err = zap.NewProduction()
	} else {
		zapLogger, err = zap.NewDevelopment()
	}

	if err != nil {
		panic("cannot initialize zap logger: " + err.Error())
	}

	Log = zapLogger.Sugar()
}

func Sync() {
	if Log != nil {
		_ = Log.Sync() // flush logs from buffer to stdout
	}
}
