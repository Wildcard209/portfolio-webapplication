package config

import (
	"os"
	"strconv"
	"time"

	"github.com/ulule/limiter/v3"
)

type EnhancedRateLimitConfig struct {
	Login   RateLimit `json:"login"`
	Refresh RateLimit `json:"refresh"`
	Upload  RateLimit `json:"upload"`
	API     RateLimit `json:"api"`
	Public  RateLimit `json:"public"`
	Admin   RateLimit `json:"admin"`
}

type RateLimit struct {
	Requests int           `json:"requests"`
	Period   time.Duration `json:"period"`
}

func (r RateLimit) ToLimiterRate() limiter.Rate {
	return limiter.Rate{
		Period: r.Period,
		Limit:  int64(r.Requests),
	}
}

func LoadRateLimitConfig() *EnhancedRateLimitConfig {
	return &EnhancedRateLimitConfig{
		Login: RateLimit{
			Requests: getEnvInt("RATE_LIMIT_LOGIN_REQUESTS", 5),
			Period:   getEnvDuration("RATE_LIMIT_LOGIN_PERIOD", "1m"),
		},
		Refresh: RateLimit{
			Requests: getEnvInt("RATE_LIMIT_REFRESH_REQUESTS", 10),
			Period:   getEnvDuration("RATE_LIMIT_REFRESH_PERIOD", "1m"),
		},
		Upload: RateLimit{
			Requests: getEnvInt("RATE_LIMIT_UPLOAD_REQUESTS", 3),
			Period:   getEnvDuration("RATE_LIMIT_UPLOAD_PERIOD", "1m"),
		},
		API: RateLimit{
			Requests: getEnvInt("RATE_LIMIT_API_REQUESTS", 60),
			Period:   getEnvDuration("RATE_LIMIT_API_PERIOD", "1m"),
		},
		Public: RateLimit{
			Requests: getEnvInt("RATE_LIMIT_PUBLIC_REQUESTS", 100),
			Period:   getEnvDuration("RATE_LIMIT_PUBLIC_PERIOD", "1m"),
		},
		Admin: RateLimit{
			Requests: getEnvInt("RATE_LIMIT_ADMIN_REQUESTS", 30),
			Period:   getEnvDuration("RATE_LIMIT_ADMIN_PERIOD", "1m"),
		},
	}
}

func getEnvInt(key string, defaultValue int) int {
	if valueStr := os.Getenv(key); valueStr != "" {
		if value, err := strconv.Atoi(valueStr); err == nil && value > 0 {
			return value
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue string) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		valueStr = defaultValue
	}

	if duration, err := time.ParseDuration(valueStr); err == nil {
		return duration
	}

	if duration, err := time.ParseDuration(defaultValue); err == nil {
		return duration
	}

	return 1 * time.Minute
}
