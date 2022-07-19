package gopool

type Config struct {
	ScaleThreshold int32
}

func NewConfig() *Config {
	return &Config{
		ScaleThreshold: 1,
	}
}
