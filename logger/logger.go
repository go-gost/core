package logger

import (
	"github.com/go-gost/core/common/util"
)

// LogFormat is format type
type LogFormat string

const (
	TextFormat LogFormat = "text"
	JSONFormat LogFormat = "json"
)

// LogLevel is Logger Level type
type LogLevel string

const (
	// DebugLevel has verbose message
	DebugLevel LogLevel = "debug"
	// InfoLevel is default log level
	InfoLevel LogLevel = "info"
	// WarnLevel is for logging messages about possible issues
	WarnLevel LogLevel = "warn"
	// ErrorLevel is for logging errors
	ErrorLevel LogLevel = "error"
	// FatalLevel is for logging fatal messages. The system shuts down after logging the message.
	FatalLevel LogLevel = "fatal"
)

type Logger interface {
	WithFields(map[string]any) Logger
	Debug(args ...any)
	Debugf(format string, args ...any)
	Info(args ...any)
	Infof(format string, args ...any)
	Warn(args ...any)
	Warnf(format string, args ...any)
	Error(args ...any)
	Errorf(format string, args ...any)
	Fatal(args ...any)
	Fatalf(format string, args ...any)
	GetLevel() LogLevel
	IsLevelEnabled(level LogLevel) bool
}

var (
	defaultLogger Logger
	defaultLoggers = make(map[string]Logger)
)

func Default() Logger {
	logger, exists := defaultLoggers[util.GetGoroutineID()]

	if exists {
		return logger
	}
	return defaultLogger
}

func SetDefault(logger Logger) {
	defaultLogger = logger
	defaultLoggers[util.GetGoroutineID()] = logger
}
