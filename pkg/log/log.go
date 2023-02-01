package log

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/kanengo/goutil/pkg/log/zaplog"
	"go.uber.org/zap"
)

var logger *zap.Logger
var logApi string

func init() {
	logger, logApi, _ = zaplog.Init()
}

func NewLogger(opts ...zaplog.ConfigOption) (*zap.Logger, string, error) {
	logger, logApi, err := zaplog.Init(opts...)
	if err != nil {
		return nil, "", err
	}

	return logger, logApi, nil
}

func InitLogger(opts ...zaplog.ConfigOption) error {
	if logger != nil {
		logger.Sync()
	}
	var err error
	logger, logApi, err = zaplog.Init(opts...)
	if err != nil {
		return err
	}

	return nil
}

func Debug(message string, fields ...zap.Field) {
	logger.Debug(message, fields...)
}

func Info(message string, fields ...zap.Field) {
	logger.Info(message, fields...)
}

func Warn(message string, fields ...zap.Field) {
	logger.Warn(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	logger.Error(message, fields...)
}

func DPanic(message string, fields ...zap.Field) {
	logger.DPanic(message, fields...)
}

func Panic(message string, fields ...zap.Field) {
	logger.Panic(message, fields...)
}

func Fatal(message string, fields ...zap.Field) {
	logger.Fatal(message, fields...)
}

func SetLogLevel(level string) error {
	level = strings.ToLower(level)
	switch strings.ToLower(level) {
	case "debug", "info", "warn", "error", "fatal":
	case "all":
		level = "debug"
	case "off", "none":
		level = "fatal"
	default:
		return errors.New("not support level")
	}
	client := http.Client{}

	type payload struct {
		Level string `json:"level"`
	}

	myPayload := payload{Level: level}
	buf, err := json.Marshal(myPayload)
	if err != nil {
		return err
	}
	Info("SetLogLevel", zap.String("path", logApi), zap.String("level", level))
	req, err := http.NewRequest("PUT", logApi, bytes.NewReader(buf))
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		Error("SetLogLevel failed", zap.Error(err), zap.String("path", logApi))
		return err
	}

	defer resp.Body.Close()

	return nil
}
