package log

import (
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestInfo(t *testing.T) {
	Info("TestInfo", zap.Time("now", time.Now()))
}

func TestError(t *testing.T) {
	Error("TestInfo", zap.Time("now", time.Now()))
}
