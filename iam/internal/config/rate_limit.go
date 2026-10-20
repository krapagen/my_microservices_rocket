package config

// rateLimitConfig — настройки distributed rate limiter поверх Redis.
// Лимитер использует redis-iam (тот же инстанс, что и сессии) — ключи
// rate:* не пересекаются с ключами сессий.
type rateLimitConfig struct {
	RatePerSecond int `yaml:"rate"  env:"RATE_LIMIT_RATE"  env-default:"100"` // лимит запросов в секунду на path
	BurstSize     int `yaml:"burst" env:"RATE_LIMIT_BURST" env-default:"200"` // размер всплеска (token bucket)
}

func (c *rateLimitConfig) Rate() int  { return c.RatePerSecond }
func (c *rateLimitConfig) Burst() int { return c.BurstSize }
