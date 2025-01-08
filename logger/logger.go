package logger

import "sync/atomic"

// LogFormat is format type
type LogFormat string

const (
	TextFormat LogFormat = "text"
	JSONFormat LogFormat = "json"
)

// LogLevel is Logger Level type
type LogLevel string

const (
	// TraceLevel has more verbose message than debug level
	TraceLevel LogLevel = "trace"
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
	Trace(args ...any)
	Tracef(format string, args ...any)
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
	defaultLogger atomic.Value
)

func Default() Logger {
	v, _ := defaultLogger.Load().(Logger)
	return v
}

func SetDefault(logger Logger) {
	defaultLogger.Store(logger)
}

type loggerGroup struct {
	loggers []Logger
}

func LoggerGroup(loggers ...Logger) Logger {
	return &loggerGroup{
		loggers: loggers,
	}
}

func (l *loggerGroup) WithFields(m map[string]any) Logger {
	lg := &loggerGroup{}
	for i := range l.loggers {
		lg.loggers = append(lg.loggers, l.loggers[i].WithFields(m))
	}
	return lg
}

func (l *loggerGroup) Trace(args ...any) {
	for _, lg := range l.loggers {
		lg.Trace(args...)
	}
}

func (l *loggerGroup) Tracef(format string, args ...any) {
	for _, lg := range l.loggers {
		lg.Tracef(format, args...)
	}
}

func (l *loggerGroup) Debug(args ...any) {
	for _, lg := range l.loggers {
		lg.Debug(args...)
	}
}

func (l *loggerGroup) Debugf(format string, args ...any) {
	for _, lg := range l.loggers {
		lg.Debugf(format, args...)
	}
}

func (l *loggerGroup) Info(args ...any) {
	for _, lg := range l.loggers {
		lg.Info(args...)
	}
}

func (l *loggerGroup) Infof(format string, args ...any) {
	for _, lg := range l.loggers {
		lg.Infof(format, args...)
	}
}

func (l *loggerGroup) Warn(args ...any) {
	for _, lg := range l.loggers {
		lg.Warn(args...)
	}
}

func (l *loggerGroup) Warnf(format string, args ...any) {
	for _, lg := range l.loggers {
		lg.Warnf(format, args...)
	}
}

func (l *loggerGroup) Error(args ...any) {
	for _, lg := range l.loggers {
		lg.Error(args...)
	}
}

func (l *loggerGroup) Errorf(format string, args ...any) {
	for _, lg := range l.loggers {
		lg.Errorf(format, args...)
	}
}

func (l *loggerGroup) Fatal(args ...any) {
	for _, lg := range l.loggers {
		lg.Fatal(args...)
	}
}

func (l *loggerGroup) Fatalf(format string, args ...any) {
	for _, lg := range l.loggers {
		lg.Fatalf(format, args...)
	}
}

func (l *loggerGroup) GetLevel() LogLevel {
	for _, lg := range l.loggers {
		return lg.GetLevel()
	}
	return InfoLevel
}

func (l *loggerGroup) IsLevelEnabled(level LogLevel) bool {
	for _, lg := range l.loggers {
		if lg.IsLevelEnabled(level) {
			return true
		}
	}
	return false
}
