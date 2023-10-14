package ipratelimit

import (
	"time"
)

type config struct {
	period time.Duration
	limit  int64
}

var defualtConfig = config{
	period: 1 * time.Hour,
	limit:  1000,
}

type Option interface {
	apply(*config)
}

type optionFunc func(*config)

func (o optionFunc) apply(c *config) {
	o(c)
}

func WithPeriod(d time.Duration) Option {
	return optionFunc(func(cfg *config) {
		if cfg.period != 0 {
			cfg.period = d
		}
	})
}

func WithLimit(v int64) Option {
	return optionFunc(func(cfg *config) {
		if cfg.limit != 0 {
			cfg.limit = v
		}
	})
}
