package ael

import (
	"context"
	"testing"

	"github.com/labstack/gommon/log"
)

func TestGetLogger_contextWithLogger(t *testing.T) {
	l := GetLogger(contextWithLogger(context.Background(), &Logger{}))
	if _, ok := l.(*Logger); !ok {
		t.Fatalf("Logger is not `*Logger`.")
	}
}

func TestGetLogger_contextBackground(t *testing.T) {
	l := GetLogger(context.Background())
	if _, ok := l.(*log.Logger); !ok {
		t.Fatalf("Logger is not `*log.Logger`.")
	}
}
