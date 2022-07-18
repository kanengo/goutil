package zaplog

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type (
	zaplogIniter interface {
		logInit(config *zaplogConfig) (zap.AtomicLevel, *zap.Logger, error)
	}

	macZaplogInit struct {
	}

	winZaplogInit struct {
	}

	unixZaplogInit struct {
	}
)

var logger *zap.Logger

func (*macZaplogInit) logInit(config *zaplogConfig) (zap.AtomicLevel, *zap.Logger, error) {
	var (
		zapconfig zap.Config
		level     zap.AtomicLevel
		logger    *zap.Logger
		err       error
	)

	//if config.isDevelop {
	//	zapconfig = zap.NewDevelopmentConfig()
	//} else {
	//	zapconfig = zap.NewProductionConfig()
	//}
	zapconfig = zap.NewDevelopmentConfig()

	zapconfig.DisableStacktrace = true
	zapconfig.EncoderConfig.TimeKey = "timestamp"
	zapconfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err = zapconfig.Build()
	level = zapconfig.Level

	return level, logger, err
}

func (*winZaplogInit) logInit(config *zaplogConfig) (zap.AtomicLevel, *zap.Logger, error) {
	var (
		zapconfig zap.Config
		level     zap.AtomicLevel
		logger    *zap.Logger
		err       error
	)

	//if config.isDevelop {
	//	zapconfig = zap.NewDevelopmentConfig()
	//} else {
	//	zapconfig = zap.NewProductionConfig()
	//}

	zapconfig = zap.NewDevelopmentConfig()

	zapconfig.DisableStacktrace = true
	zapconfig.EncoderConfig.TimeKey = "timestamp"
	zapconfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err = zapconfig.Build()
	level = zapconfig.Level

	return level, logger, err
}

func epochFullTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func (*unixZaplogInit) logInit(config *zaplogConfig) (zap.AtomicLevel, *zap.Logger, error) {
	var (
		level  zap.AtomicLevel
		logger *zap.Logger
		err    error
	)

	writers := []zapcore.WriteSyncer{os.Stderr}
	output := zapcore.NewMultiWriteSyncer(writers...)
	if config.logPath != "" {
		output = zapcore.AddSync(&lumberjack.Logger{
			Filename:   config.logPath,
			MaxSize:    config.maxSize,
			MaxAge:     config.maxAge,
			Compress:   config.compress,
			MaxBackups: config.maxBackups,
			LocalTime:  true,
		})
	}
	encConfig := zap.NewProductionEncoderConfig()
	encConfig.TimeKey = "timestamp"
	encConfig.EncodeTime = epochFullTimeEncoder

	encoder := zapcore.NewJSONEncoder(encConfig)
	if config.isDevelop {
		level = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else {
		level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	logger = zap.New(zapcore.NewCore(encoder, output, level), zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return level, logger, err
}

func Init(opts ...ConfigOption) (*zap.Logger, string, error) {
	var (
		logIniter zaplogIniter
		level     zap.AtomicLevel
		logger    *zap.Logger
		err       error
	)

	config := defaultConfig
	for _, opt := range opts {
		opt.apply(&config)
	}

	if runtime.GOOS == "darwin" {
		logIniter = &macZaplogInit{}
	} else if runtime.GOOS == "windows" {
		logIniter = &winZaplogInit{}
	} else {
		logIniter = &unixZaplogInit{}
	}

	if level, logger, err = logIniter.logInit(&config); err != nil {
		fmt.Printf("zaplogInit failed, err:%v\n", err)
		return logger, "", err
	}
	//
	//if config.withPid {
	//	llog = llog.With(zap.Int("pid", os.Getpid()))
	//}
	//

	if len(config.fields) > 0 {
		for key, value := range config.fields {
			logger = logger.With(zap.String(key, value))
		}
	}

	logger = logger.WithOptions(zap.AddCallerSkip(1))

	logApi := LogLevelHttpServer(&config, level)

	return logger, logApi, nil
}
