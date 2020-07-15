package logging

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"time"

	"cloud.google.com/go/logging"
	"github.com/labstack/gommon/log"
	logpb "google.golang.org/genproto/googleapis/logging/v2"
)

var severityLogLevel = map[logging.Severity]log.Lvl{
	logging.Default:  0,
	logging.Debug:    log.DEBUG,
	logging.Info:     log.INFO,
	logging.Warning:  log.WARN,
	logging.Error:    log.ERROR,
	logging.Critical: 6,
	logging.Alert:    7,
}

// Logger is echo logger using cloud logging.
type Logger struct {
	logger *logging.Logger

	trace  string
	spanID string
	level  log.Lvl

	maxSeverity logging.Severity
}

// NewLogger returns echo logger using cloud logging.
func NewLogger(logger *logging.Logger, trace, spanID string) *Logger {
	return &Logger{
		logger: logger,
		trace:  trace,
		spanID: spanID,
	}
}

// Output does nothing.
func (l *Logger) Output() io.Writer { return nil }

// SetOutput does nothing.
func (l *Logger) SetOutput(w io.Writer) {}

// Prefix does nothing.
func (l *Logger) Prefix() string { return "" }

// SetPrefix does nothing.
func (l *Logger) SetPrefix(p string) {}

// SetHeader does nothing.
func (l *Logger) SetHeader(h string) {}

// Level returns print log level.
func (l *Logger) Level() log.Lvl {
	return l.level
}

// SetLevel sets print log level.
func (l *Logger) SetLevel(v log.Lvl) {
	l.level = v
}

// Print prints i.
func (l *Logger) Print(i ...interface{}) {
	l.log(logging.Default, fmt.Sprint(i...))
}

// Printf prints format string.
func (l *Logger) Printf(format string, args ...interface{}) {
	l.log(logging.Default, fmt.Sprintf(format, args...))
}

// Printj prints j.
func (l *Logger) Printj(j log.JSON) {
	l.log(logging.Default, j)
}

// Debug prints i to log level debug.
func (l *Logger) Debug(i ...interface{}) {
	l.log(logging.Debug, fmt.Sprint(i...))
}

// Debugf prints format string to log level debug.
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.log(logging.Debug, fmt.Sprintf(format, args...))
}

// Debugj prints j to log level debug.
func (l *Logger) Debugj(j log.JSON) {
	l.log(logging.Debug, j)
}

// Info prints i to log level info.
func (l *Logger) Info(i ...interface{}) {
	l.log(logging.Info, fmt.Sprint(i...))
}

// Infof prints format string to log level info.
func (l *Logger) Infof(format string, args ...interface{}) {
	l.log(logging.Info, fmt.Sprintf(format, args...))
}

// Infoj prints j to log level info.
func (l *Logger) Infoj(j log.JSON) {
	l.log(logging.Info, j)
}

// Warn prints i to log level warning.
func (l *Logger) Warn(i ...interface{}) {
	l.log(logging.Warning, fmt.Sprint(i...))
}

// Warnf prints format string to log level warning.
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.log(logging.Warning, fmt.Sprintf(format, args...))
}

// Warnj prints j to log level warning.
func (l *Logger) Warnj(j log.JSON) {
	l.log(logging.Warning, j)
}

// Error prints i to log level error.
func (l *Logger) Error(i ...interface{}) {
	l.log(logging.Error, fmt.Sprint(i...))
}

// Errorf prints format string to log level error.
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.log(logging.Error, fmt.Sprintf(format, args...))
}

// Errorj prints j to log level error.
func (l *Logger) Errorj(j log.JSON) {
	l.log(logging.Error, j)
}

// Fatal prints i and exit 1.
func (l *Logger) Fatal(i ...interface{}) {
	l.log(logging.Critical, fmt.Sprint(i...))
	l.logger.Flush()
	os.Exit(1)
}

// Fatalf prints format string and exit 1.
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.log(logging.Critical, fmt.Sprintf(format, args...))
	l.logger.Flush()
	os.Exit(1)
}

// Fatalj prints j and exit 1.
func (l *Logger) Fatalj(j log.JSON) {
	l.log(logging.Critical, j)
	l.logger.Flush()
	os.Exit(1)
}

// Panic prints i and panic.
func (l *Logger) Panic(i ...interface{}) {
	s := fmt.Sprint(i...)
	l.log(logging.Alert, s)
	panic(s)
}

// Panicf prints format string and panic.
func (l *Logger) Panicf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	l.log(logging.Alert, s)
	panic(s)
}

// Panicj prints j and panic.
func (l *Logger) Panicj(j log.JSON) {
	l.log(logging.Alert, j)
	panic(j)
}

func (l *Logger) log(severity logging.Severity, payload interface{}) {
	if l.level > severityLogLevel[severity] {
		return
	}

	if severity > l.maxSeverity {
		l.maxSeverity = severity
	}

	pc, file, line, _ := runtime.Caller(2)
	f := runtime.FuncForPC(pc)
	l.logger.Log(logging.Entry{
		Timestamp:    time.Now(),
		Severity:     severity,
		Payload:      payload,
		Trace:        l.trace,
		SpanID:       l.spanID,
		TraceSampled: true,
		SourceLocation: &logpb.LogEntrySourceLocation{
			File:     file,
			Line:     int64(line),
			Function: f.Name(),
		},
	})
}
