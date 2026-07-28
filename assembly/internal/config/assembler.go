package config

type assemblerConfig struct {
	AssembleLimitTimeSec int64 `yaml:"limit" env-default:"5"`
}

func (c *assemblerConfig) AssembleLimitTimeLimit() int64 {
	return c.AssembleLimitTimeSec
}
