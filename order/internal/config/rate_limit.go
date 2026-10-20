package config

// internal/config/rate_limit.go
// rateLimitConfig — внутренняя структура пакета config (имя с маленькой буквы),
// наружу торчит только метод RateLimit() RateLimitConfig интерфейса.
type rateLimitConfig struct {
	RedisAddr     string `yaml:"redis_address" env:"RATE_LIMIT_REDIS_ADDRESS"` // host:port redis-ratelimit
	RatePerSecond int    `yaml:"rate"          env:"RATE_LIMIT_RATE"`          // лимит запросов в секунду на path
	BurstSize     int    `yaml:"burst"         env:"RATE_LIMIT_BURST"`         // размер всплеска (token bucket)
}

// RateLimitConfig — публичный интерфейс настроек distributed rate limiter.
type RateLimitConfig interface {
	// RedisAddress — адрес Redis для rate limiter (host:port).
	RedisAddress() string
	// Rate — лимит запросов в секунду на один path.
	Rate() int
	// Burst — максимальный всплеск запросов.
	Burst() int
}

func (c *rateLimitConfig) RedisAddress() string { return c.RedisAddr }
func (c *rateLimitConfig) Rate() int            { return c.RatePerSecond }
func (c *rateLimitConfig) Burst() int           { return c.BurstSize }

// RateLimit возвращает конфигурацию rate limiter.
// Поле названо RateLimitCfg, чтобы не конфликтовать по имени с этим методом.
func (c *config) RateLimit() RateLimitConfig {
	return &c.RateLimitCfg
}
