// Package logger defines the structured logging interface and a global
// default logger that components use when no explicit logger is configured.
package logger

import "sync/atomic"

// LogFormat is the log output format type.
type LogFormat string

const (
	// TextFormat produces human-readable log output.
	TextFormat LogFormat = "text"
	// JSONFormat produces structured JSON log output.
	JSONFormat LogFormat = "json"
)

// LogLevel is the severity level of a log entry.
type LogLevel string

const (
	// TraceLevel has more verbose messages than DebugLevel.
	TraceLevel LogLevel = "trace"
	// DebugLevel is for verbose diagnostic messages.
	DebugLevel LogLevel = "debug"
	// InfoLevel is the default log level for informational messages.
	InfoLevel LogLevel = "info"
	// WarnLevel is for messages about potential issues.
	WarnLevel LogLevel = "warn"
	// ErrorLevel is for error messages indicating a failure.
	ErrorLevel LogLevel = "error"
	// FatalLevel is for fatal errors. The program typically exits after logging.
	FatalLevel LogLevel = "fatal"
)

// Logger is the structured logging interface. Implementations may write to
// files, stdout, sockets, or external log services.
type Logger interface {
	// WithFields returns a new Logger with the given fields attached to
	// every subsequent log entry.
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
	// GetLevel returns the current log level.
	GetLevel() LogLevel
	// IsLevelEnabled reports whether messages at the given level are logged.
	IsLevelEnabled(level LogLevel) bool
}

var (
	// defaultLogger holds the global default Logger, stored atomically.
	defaultLogger atomic.Value
)

// Default returns the global default Logger. It may return nil if no default
// has been set.
func Default() Logger {
	v, _ := defaultLogger.Load().(Logger)
	return v
}

// SetDefault sets the global default Logger.
func SetDefault(logger Logger) {
	defaultLogger.Store(logger)
}

// loggerGroup is a Logger that fans out messages to multiple underlying loggers.
type loggerGroup struct {
	loggers []Logger
}

// LoggerGroup creates a Logger that broadcasts all log calls to the given
// set of loggers. This is useful for emitting logs to multiple destinations
// simultaneously (e.g. both a file and stdout).
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
