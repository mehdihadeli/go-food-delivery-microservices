package redis

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

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

func NewUniversalRedisClient(cfg *RedisOptions) redis.UniversalClient {
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
