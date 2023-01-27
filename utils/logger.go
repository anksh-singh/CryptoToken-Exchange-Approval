package utils

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"net/url"
)

var sugaredLogger *zap.SugaredLogger

type lumberjackSink struct {
	*lumberjack.Logger
}

func (lumberjackSink) Sync() error {
	return nil
}

func SetupLogger(loglevel,logFile,encodingFormat string) *zap.SugaredLogger {
	writer := lumberjack.Logger{
			Filename:  logFile,
			MaxSize:    50, //megabytes
			MaxBackups: 3,
			MaxAge:  28, //days
			Compress: true,
		}
	err := zap.RegisterSink("lumberjack", func(*url.URL) (zap.Sink, error) {
		return lumberjackSink{
			Logger: &writer,
		}, nil
	})
	if err != nil {
		panic(fmt.Sprintf("xap register sink error: %v", err))
		return nil
	}
	cfg := zap.Config{
		Encoding:    encodingFormat,
		Level:       zap.NewAtomicLevelAt(getLogLevel(loglevel)),
		OutputPaths: []string{"stderr",fmt.Sprintf("lumberjack:%s", logFile)},
		ErrorOutputPaths: []string{"stderr",fmt.Sprintf("lumberjack:%s", logFile)},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,

			MessageKey: "msg",
		}}
	_logger, err  := cfg.Build()
	if err != nil {
		panic(fmt.Sprintf("build zap logger from config error: %v", err))
	}
	zap.ReplaceGlobals(_logger)
	sugaredLogger = _logger.Named("").Sugar()
	return sugaredLogger
}

func getLogLevel(level string) zapcore.Level {
	switch level {
	case "DEBUG", "debug":
		return zapcore.DebugLevel
	case "ERROR", "error":
		return zapcore.ErrorLevel
	case "WARN", "warn":
		return zapcore.WarnLevel
	case "INFO", "info":
		return zapcore.InfoLevel
	default:
		return zapcore.ErrorLevel
	}
}