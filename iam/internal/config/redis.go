package config

import (
	"fmt"
	"time"
)

type redisConfig struct {
	Host           string        `yaml:"host"            env:"REDIS_HOST"     env-default:"localhost"`
	Port           string        `yaml:"port"            env:"REDIS_PORT"     env-default:"6379"`
	Password       string        `yaml:"password"        env:"REDIS_PASSWORD" env-default:""`
	DB             int           `yaml:"db"              env:"REDIS_DB"       env-default:"0"`
	ConnectTimeout time.Duration `yaml:"connect_timeout" env:"REDIS_CT"       env-default:"5s"`
	MaxIdle        int           `yaml:"max_idle"        env:"REDIS_MAX_IDLE" env-default:"10"`
	IdleTimeout    time.Duration `yaml:"idle_timeout"    env:"REDIS_IDLE_TIMEOUT" env-default:"10s"`
}

func (r *redisConfig) Address() string {
	return fmt.Sprintf("%s:%s", r.Host, r.Port)
}

func (r *redisConfig) Database() int {
	return r.DB
}

func (r *redisConfig) Pass() string {
	return r.Password
}

func (r *redisConfig) ConT() time.Duration {
	return r.ConnectTimeout
}

func (r *redisConfig) MIdle() int {
	return r.MaxIdle
}

func (r *redisConfig) IdleT() time.Duration {
	return r.IdleTimeout
}
