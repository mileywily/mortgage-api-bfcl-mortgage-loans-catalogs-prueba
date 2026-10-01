package loggers

import (
	"context"
	"io"
	"log/slog"
	"os"

	loggerFif "github.com/falabella-regulado/go-lib-logger-fif"
)

// Slog DEFAULT LEVEL IS INFO
type Slog struct {
	logger *slog.Logger
}

var logLevelMap = map[loggerFif.LogLevel]slog.Level{
	loggerFif.DebugLevel: slog.LevelDebug,
	loggerFif.InfoLevel:  slog.LevelInfo,
	loggerFif.WarnLevel:  slog.LevelWarn,
	loggerFif.ErrorLevel: slog.LevelError,
}

func getSlogLevel(level loggerFif.LogLevel) slog.Level {
	newLevel, isOk := logLevelMap[level]
	if !isOk {
		return slog.LevelInfo
	}
	return newLevel
}

func NewJsonSlogLogger(level loggerFif.LogLevel, writer io.Writer) loggerFif.Logger {
	slogLevel := getSlogLevel(level)
	if writer == nil {
		writer = os.Stdout
	}
	logger := slog.New(slog.NewJSONHandler(writer, &slog.HandlerOptions{Level: slogLevel}))
	return &Slog{
		logger: logger,
	}
}

func NewTextSlogLogger(level loggerFif.LogLevel, writer io.Writer) loggerFif.Logger {
	slogLevel := getSlogLevel(level)
	if writer == nil {
		writer = os.Stdout
	}
	logger := slog.New(slog.NewTextHandler(writer, &slog.HandlerOptions{Level: slogLevel}))
	return &Slog{
		logger: logger,
	}
}

func NewNoopSlogLogger() loggerFif.Logger {
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo}))
	return &Slog{
		logger: logger,
	}
}

func (s Slog) Log(level loggerFif.LogLevel, msg string, args ...interface{}) {
	slogLevel := getSlogLevel(level)
	s.logger.Log(context.Background(), slogLevel, msg, args...)
}

func (s Slog) Debug(msg string, args ...interface{}) {
	s.Log(loggerFif.DebugLevel, msg, args...)
}

func (s Slog) Info(msg string, args ...interface{}) {
	s.Log(loggerFif.InfoLevel, msg, args...)
}

func (s Slog) Error(msg string, args ...interface{}) {
	s.Log(loggerFif.ErrorLevel, msg, args...)
}

func (s Slog) Warn(msg string, args ...interface{}) {
	s.Log(loggerFif.WarnLevel, msg, args...)
}
