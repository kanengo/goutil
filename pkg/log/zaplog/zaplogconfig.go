package zaplog

const (
	LOGPATH = ""
)

type (
	zaplogConfig struct {
		logPath    string
		isDevelop  bool
		maxSize    int
		maxAge     int
		maxBackups int
		compress   bool
		logApiPath string
		listenAddr string
		fields     map[string]string
	}

	ConfigOption interface {
		apply(config *zaplogConfig)
	}
)

var defaultConfig = zaplogConfig{
	logPath:    LOGPATH,
	isDevelop:  false,
	maxSize:    100,
	maxAge:     3,
	compress:   true,
	logApiPath: "/log",
	listenAddr: "0.0.0.0:0",
}

type zaplogOptionFunc func(config *zaplogConfig)

func (f zaplogOptionFunc) apply(config *zaplogConfig) {
	f(config)
}

func SetLogPath(path string) ConfigOption {
	return zaplogOptionFunc(func(config *zaplogConfig) {
		config.logPath = path
	})

}

func SetIsDevelop(flag bool) ConfigOption {
	return zaplogOptionFunc(func(config *zaplogConfig) {
		config.isDevelop = flag
	})
}

func SetMaxSize(size int) ConfigOption {
	return zaplogOptionFunc(func(config *zaplogConfig) {
		config.maxSize = size
	})
}

func SetMaxAge(age int) ConfigOption {
	return zaplogOptionFunc(func(config *zaplogConfig) {
		config.maxAge = age
	})
}

func SetMaxBackups(backups int) ConfigOption {
	return zaplogOptionFunc(func(config *zaplogConfig) {
		config.maxBackups = backups
	})
}

func SetCompress(compress bool) ConfigOption {
	return zaplogOptionFunc(func(config *zaplogConfig) {
		config.compress = compress
	})
}

func SetLogApiPath(logApiPath string) ConfigOption {
	return zaplogOptionFunc(func(config *zaplogConfig) {
		config.logApiPath = logApiPath
	})
}

func SetListenAddr(listenAddr string) ConfigOption {
	return zaplogOptionFunc(func(config *zaplogConfig) {
		config.listenAddr = listenAddr
	})
}

func SetFiled(key, value string) ConfigOption {
	return zaplogOptionFunc(func(config *zaplogConfig) {
		if config.fields == nil {
			config.fields = make(map[string]string)
		}
		config.fields[key] = value
	})
}
