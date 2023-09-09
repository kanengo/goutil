package log

import (
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestInfo(t *testing.T) {
	InitLogger()
	logger.Info("TestInfo", zap.Time("now", time.Now()))
	SetLogLevel("debug")
	logger.Debug("debug")
}

func TestError(t *testing.T) {
	Error("TestInfo", zap.Time("now", time.Now()))
}

func TestBeforeLogHandler(t *testing.T) {
	SetBeforeLogHandler(func() []zap.Field {
		return []zap.Field{zap.Int64("uid", 6544)}
	})
	Info("TestBeforeLogHandler", zap.Time("now", time.Now()))
	Debug("TestBeforeLogHandler", zap.Time("now", time.Now()))
}
