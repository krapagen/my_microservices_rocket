package config

type loggerConfig struct {
	Level string `yaml:"level" env:"LOGGER_LEVEL" env-default:"info"`
}

func (c *loggerConfig) LoggerLevel() string {
	return "level=" + c.Level
}
