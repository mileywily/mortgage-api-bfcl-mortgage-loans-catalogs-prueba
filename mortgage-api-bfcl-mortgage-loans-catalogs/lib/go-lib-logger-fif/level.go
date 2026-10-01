package logger_fif

import (
	"errors"
	"strings"
)

type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	Debug = "DEBUG"
	Info  = "INFO"
	Warn  = "WARN"
	Error = "ERROR"
)

func StringToLogLevel(level string) (LogLevel, error) {
	switch strings.ToUpper(level) {
	case Debug:
		return DebugLevel, nil
	case Info:
		return InfoLevel, nil
	case Warn:
		return WarnLevel, nil
	case Error:
		return ErrorLevel, nil
	default:
		return 0, errors.New("invalid log level")
	}
}
