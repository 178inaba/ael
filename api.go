package log

import (
	"context"

	"cloud.google.com/go/logging"
)

type Logger struct {
	logger *logging.Logger
}

func New(ctx context.Context, parent string) (*Logger, error) {
	c, err := logging.NewClient(ctx, parent)
	if err != nil {
		return nil, err
	}

	return &Logger{logger: c.Logger("")}, nil
}

// Debugf formats its arguments according to the format, analogous to fmt.Printf,
// and records the text as a log message at Debug level. The message will be associated
// with the request linked with the provided context.
func (l *Logger) Debugf(ctx context.Context, format string, args ...interface{}) {
}

// Infof is like Debugf, but at Info level.
func (l *Logger) Infof(ctx context.Context, format string, args ...interface{}) {
}

// Warningf is like Debugf, but at Warning level.
func (l *Logger) Warningf(ctx context.Context, format string, args ...interface{}) {
}

// Errorf is like Debugf, but at Error level.
func (l *Logger) Errorf(ctx context.Context, format string, args ...interface{}) {
}

// Criticalf is like Debugf, but at Critical level.
func (l *Logger) Criticalf(ctx context.Context, format string, args ...interface{}) {
}
