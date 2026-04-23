package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger
var L *zap.Logger

func InitLogger() {
	// App Env
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	// Encoder config
	config := zap.NewProductionEncoderConfig()
	
	// Formating logging time
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	
	// Formating logging level
	config.EncodeLevel = zapcore.CapitalColorLevelEncoder

	// Formating logging caller
	config.EncodeCaller = zapcore.ShortCallerEncoder
	
	// Formating logging
	var encoder zapcore.Encoder
	if env == "development" {
		encoder = zapcore.NewConsoleEncoder(config)
	} else {
		encoder = zapcore.NewJSONEncoder(config)
	}

	core := zapcore.NewCore(
		encoder, 
		zapcore.AddSync(os.Stdout), 
		zap.InfoLevel,
	)

	// Add caller for trigger logging
	Log = zap.New(core, zap.AddCaller())
	L = Log
}