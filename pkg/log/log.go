// Package log provides a globally accessible zap logger which is initialized at startup.
package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	globalLogger *zap.Logger
)

const (
	DebugLevelStr   = "debug"
	InfoLevelStr    = "info"
	WarningLevelStr = "warn"
	ErrorLevelStr   = "error"
)

// Initialize sets up the logger with provided log level and log rolling with lumberjack. Log level defaults to "warn".
//
// - logLevel: Level of logs to display can be debig, info, warn or error use log.*LevelStr constants.
// - disableConsoleOutput: to disable output to console.
func Initialize(logLevel string, disableConsoleOutput bool) {
	var level zapcore.Level

	switch logLevel {
	case DebugLevelStr:
		level = zap.DebugLevel
	case InfoLevelStr:
		level = zap.InfoLevel
	case WarningLevelStr:
		level = zap.WarnLevel
	case ErrorLevelStr:
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}

	lumberjackLogger := &lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    10, // MB
		MaxBackups: 5,
		MaxAge:     10,   // Days
		Compress:   true, // .gz
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("01/02/2006 15:04:05"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	writeSyncer := zapcore.AddSync(lumberjackLogger)

	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		writeSyncer,
		level,
	)

	cores := []zapcore.Core{fileCore}

	if !disableConsoleOutput {
		consoleCore := zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			level,
		)

		cores = append(cores, consoleCore)
	}

	combinedCore := zapcore.NewTee(cores...)

	logger := zap.New(combinedCore, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	globalLogger = logger
}

// Logger returns the initialized logger.
func Logger() *zap.Logger {
	return globalLogger
}
