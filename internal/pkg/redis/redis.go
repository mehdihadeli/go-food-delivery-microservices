package redis

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Database int    `mapstructure:"database"`
	PoolSize int    `mapstructure:"poolSize"`
}

const (
	maxRetries      = 5
	minRetryBackoff = 300 * time.Millisecond
	maxRetryBackoff = 500 * time.Millisecond
	dialTimeout     = 5 * time.Second
	readTimeout     = 5 * time.Second
	writeTimeout    = 3 * time.Second
	minIdleConns    = 20
	poolTimeout     = 6 * time.Second
	idleTimeout     = 12 * time.Second
)

func NewUniversalRedisClient(cfg *RedisConfig) redis.UniversalClient {
	universalClient := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:           []string{fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)},
		Password:        cfg.Password, // no password set
		DB:              cfg.Database, // use defaultLogger Database
		MaxRetries:      maxRetries,
		MinRetryBackoff: minRetryBackoff,
		MaxRetryBackoff: maxRetryBackoff,
		DialTimeout:     dialTimeout,
		ReadTimeout:     readTimeout,
		WriteTimeout:    writeTimeout,
		PoolSize:        cfg.PoolSize,
		MinIdleConns:    minIdleConns,
		PoolTimeout:     poolTimeout,
		IdleTimeout:     idleTimeout,
	})

	return universalClient
}
