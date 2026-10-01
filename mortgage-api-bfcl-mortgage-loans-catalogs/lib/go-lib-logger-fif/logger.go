package logger_fif

type Logger interface {
	Log(LogLevel, string, ...interface{})
	Debug(string, ...interface{})
	Info(string, ...interface{})
	Error(string, ...interface{})
	Warn(string, ...interface{})
}
